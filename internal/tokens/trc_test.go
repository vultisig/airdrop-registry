package tokens

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func TestTrcAutoDiscovery(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   trcAccountResponse
		expectedCoins  []models.Coin
		expectedError  bool
	}{
		{
			name:           "successful response with TRX balance",
			responseStatus: http.StatusOK,
			responseBody: trcAccountResponse{
				Success: true,
				Data: []trcAccount{
					{
						Balance: 1000000,
						Trc20:   []map[string]string{},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create service with test server URL
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()
			// Create test server
			cmcService, err := NewCMCService()
			if err != nil {
				t.Fatalf("Failed to create CMCService: %v", err)
			}
			trc := &trcDiscoveryService{
				logger:      logrus.New(),
				tronBaseURL: server.URL,
				cmcService:  cmcService,
			}
			trc.cmcService.baseURL = "https://api.vultisig.com/cmc/v1/cryptocurrency"

			// Call method
			coins, err := trc.Discover("TS98H6jSx6uv1gG1vx6CJZMeYGkMZXgQ7K", common.Tron)

			// Check results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedCoins == nil {
					assert.Nil(t, coins)
				} else {
					for i := range tt.expectedCoins {
						for _, coin := range coins {
							if coin.ContractAddress == tt.expectedCoins[i].ContractAddress {
								assert.Equal(t, tt.expectedCoins[i].ContractAddress, coin.ContractAddress)
								assert.Equal(t, tt.expectedCoins[i].Decimals, coin.Decimals)
							}
						}
					}
				}
			}
		})
	}
}

//go:embed trc_symbol.json
var decimalResponse string

//go:embed trc_decimal.json
var symbolResponse string

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
