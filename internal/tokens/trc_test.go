package tokens

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func TestTrcAutoDiscovery(t *testing.T) {
	test := struct {
		name           string
		responseStatus int
		responseBody   trcAccountResponse
		expectedCMCID  int
		expectedCoins  []models.Coin
		expectedError  bool
	}{
		name:           "successful response with TRX balance",
		responseStatus: http.StatusOK,
		responseBody: trcAccountResponse{
			Success: true,
			Data: []trcAccount{
				{
					Balance: 1000000,
					Trc20: []map[string]string{
						{
							"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t": "5948041202973",
						},
					},
				},
			},
		},
		expectedCoins: []models.Coin{
			{
				ContractAddress: "USDT",
				Decimals:        18,
			},
		},
		expectedError: false,
	}
	// Create service with test server URL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(test.responseStatus)
		json.NewEncoder(w).Encode(test.responseBody)
	}))
	defer server.Close()
	// Create test server
	cmcService := CMCService{
		logger:        logrus.New(),
		baseURL:       server.URL,
		cachedData:    cache.New(10*time.Hour, 1*time.Hour),
		nativeCoinIds: map[string]int{},
	}
	trc := &trcDiscoveryService{
		logger:      logrus.New(),
		tronBaseURL: server.URL,
		cmcService:  &cmcService,
	}
	trc.cmcService.cachedData.Set(trc.cmcService.getCacheKey("Tron", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"), 825, cache.DefaultExpiration)

	// Call method
	coins, err := trc.Discover("TS98H6jSx6uv1gG1vx6CJZMeYGkMZXgQ7K", common.Tron)

	// Check results
	if test.expectedError {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		if test.expectedCoins == nil {
			assert.Nil(t, coins)
		} else {
			for i := range test.expectedCoins {
				for _, coin := range coins {
					if coin.ContractAddress == test.expectedCoins[i].ContractAddress {
						assert.Equal(t, test.expectedCoins[i].ContractAddress, coin.ContractAddress)
						assert.Equal(t, test.expectedCoins[i].Decimals, coin.Decimals)
					}
				}
			}
		}
	}
}

//go:embed trc_symbol.json
var symbolResponse string

//go:embed trc_decimal.json
var decimalResponse string

func TestFetchTokenData(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedSymbol string
		expectedError  bool
	}{
		{
			name:           "successful response with TRX symbol",
			responseStatus: http.StatusOK,
			responseBody:   symbolResponse,
			expectedSymbol: "USDT",
			expectedError:  false,
		},
		{
			name:           "successful response with TRX symbol",
			responseStatus: http.StatusOK,
			responseBody:   decimalResponse,
			expectedSymbol: "USDT",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create service with test server URL
			cmcService, err := NewCMCService()
			if err != nil {
				t.Fatalf("Failed to create CMCService: %v", err)
			}
			td := &trcDiscoveryService{
				logger:      logrus.New(),
				tronBaseURL: "https://api.trongrid.io",
				cmcService:  cmcService,
			}

			// Call method
			address := "TS98H6jSx6uv1gG1vx6CJZMeYGkMZXgQ7K"
			contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
			selector := "symbol()"
			symbol, err := td.fetchTokenData(address, contract, selector)

			// Check results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSymbol, symbol)
			}
		})
	}

}

func Test_trcDiscoveryService_Search(t *testing.T) {
	// Setup test server with both symbol and decimal responses
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/wallet/triggerconstantcontract" {
			// Check the function selector in the request body to determine which response to send
			var requestBody struct {
				FunctionSelector string `json:"function_selector"`
			}
			json.NewDecoder(r.Body).Decode(&requestBody)

			switch requestBody.FunctionSelector {
			case "symbol()":
				w.Write([]byte(symbolResponse))
			case "decimals()":
				w.Write([]byte(decimalResponse))
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}
		// Default response for other endpoints
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(trcAccountResponse{
			Success: true,
			Data: []trcAccount{
				{
					Balance: 1000000,
					Trc20:   []map[string]string{},
				},
			},
		})
	}))
	defer mockServer.Close()

	// Initialize discovery service
	discovery := setupTrcDiscoveryService(mockServer.URL)

	// Test cases
	testCases := []struct {
		name     string
		input    models.CoinBase
		expected models.CoinBase
	}{
		{
			name: "USDT token search",
			input: models.CoinBase{
				Ticker:          "USDT",
				Address:         "TS98H6jSx6uv1gG1vx6CJZMeYGkMZXgQ7K",
				ContractAddress: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
				Decimals:        6,
				Chain:           common.Tron,
			},
			expected: models.CoinBase{
				Ticker:          "USDT",
				Address:         "TS98H6jSx6uv1gG1vx6CJZMeYGkMZXgQ7K",
				ContractAddress: "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
				Decimals:        6,
				Chain:           common.Tron,
				CMCId:           825,
			},
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := discovery.Search(tc.input)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func setupTrcDiscoveryService(baseURL string) *trcDiscoveryService {
	cmcService, _ := NewCMCService()
	discovery := &trcDiscoveryService{
		logger:      logrus.New(),
		tronBaseURL: baseURL,
		cmcService:  cmcService,
	}

	// Setup cache
	discovery.cmcService.cachedData.Set(
		discovery.cmcService.getCacheKey("Tron", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"),
		825,
		cache.DefaultExpiration,
	)

	return discovery
}
