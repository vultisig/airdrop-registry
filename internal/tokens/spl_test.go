package tokens

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

//go:embed sol_token_response.json
var sol_token_response string

//go:embed sol_token_info.json
var sol_token_info string

func TestSolDiscoveryService_processTokenAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sol_token_response))
	}))
	CMCService := CMCService{
		cachedData: cache.New(5*time.Minute, 10*time.Minute),
		baseURL:    server.URL,
		logger:     logrus.New(),
	}

	CMCService.cachedData.Set(CMCService.getCacheKey("Solana", "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"), 3408, cache.DefaultExpiration)
	service := &splDiscoveryService{
		logger:      logrus.New(),
		cmcService:  &CMCService,
		baseAddress: server.URL,
	}

	// Parse the embedded JSON response
	var response jsonRpcResponse
	err := json.Unmarshal([]byte(sol_token_response), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal embedded JSON: %v", err)
	}

	testAddr := "CG4V2eoUXnwJSDsmr1fNdbR9r63XHLKD9gA2xpCRdRby"
	results, err := service.fetchTokenAccounts(testAddr)

	// Assertions
	if err != nil {
		t.Errorf("processTokenAccounts() error = %v", err)
		return
	}

	if len(results) == 0 {
		t.Error("processTokenAccounts() returned empty results")
		return
	}

	expected := models.CoinBase{
		Address:         testAddr,
		Balance:         "24389303",
		Chain:           common.Solana,
		IsNative:        false,
		ContractAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		CMCId:           3408,
	}
	if len(results) == 0 {
		t.Error("processTokenAccounts() returned empty results")
		return
	}

	// Validate that all required fields in the results are populated
	for i, result := range results {
		if result.Address == "" {
			t.Errorf("processTokenAccounts() result[%d] has an empty Address field", i)
		}
		if result.Balance == "" {
			t.Errorf("processTokenAccounts() result[%d] has an empty Balance field", i)
		}
		if result.Chain == 0 {
			t.Errorf("processTokenAccounts() result[%d] has an uninitialized Chain field", i)
		}
		if result.ContractAddress == "" {
			t.Errorf("processTokenAccounts() result[%d] has an empty ContractAddress field", i)
		}
		if result.CMCId == 0 {
			t.Errorf("processTokenAccounts() result[%d] has an uninitialized CMCId field", i)
		}
	}

	if !reflect.DeepEqual(results[0].Address, expected.Address) ||
		!reflect.DeepEqual(results[0].Balance, expected.Balance) ||
		!reflect.DeepEqual(results[0].Chain, expected.Chain) ||
		!reflect.DeepEqual(results[0].IsNative, expected.IsNative) ||
		!reflect.DeepEqual(results[0].ContractAddress, expected.ContractAddress) ||
		!reflect.DeepEqual(results[0].CMCId, expected.CMCId) {
		t.Errorf("processTokenAccounts() = %v, want %v", results[0], expected)
	}
}

func TestSolDiscoveryService_Search(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(sol_token_info))
	}))
	defer server.Close()

	// Initialize discovery service
	discovery := setupSolDiscoveryService(server.URL)

	// Test cases
	testCases := []struct {
		name     string
		input    models.CoinBase
		expected models.CoinBase
		wantErr  bool
	}{
		{
			name: "USDC token search",
			input: models.CoinBase{
				Ticker:          "USDC",
				Address:         "CG4V2eoUXnwJSDsmr1fNdbR9r63XHLKD9gA2xpCRdRby",
				ContractAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				Chain:           common.Solana,
				Decimals:        6,
			},
			wantErr: false,
			expected: models.CoinBase{
				Ticker:          "USDC",
				Address:         "CG4V2eoUXnwJSDsmr1fNdbR9r63XHLKD9gA2xpCRdRby",
				ContractAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				Chain:           common.Solana,
				CMCId:           3408,
				Decimals:        6,
			},
		},
		{
			name: "Invalid address search",
			input: models.CoinBase{
				Ticker:          "INVALID",
				Address:         "invalid_address",
				ContractAddress: "invalid_contract",
				Chain:           common.Solana,
			},
			wantErr:  true,
			expected: models.CoinBase{},
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := discovery.Search(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("Search() error = %v", err)
				return
			}
			if !tc.wantErr && !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Search() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func setupSolDiscoveryService(baseURL string) *splDiscoveryService {
	cmcService := &CMCService{
		cachedData: cache.New(5*time.Minute, 10*time.Minute),
		baseURL:    baseURL,
		logger:     logrus.New(),
	}

	// Setup cache
	cmcService.cachedData.Set(
		cmcService.getCacheKey("Solana", "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"),
		3408,
		cache.DefaultExpiration,
	)

	return &splDiscoveryService{
		logger:      logrus.New(),
		cmcService:  cmcService,
		baseAddress: baseURL,
	}
}
