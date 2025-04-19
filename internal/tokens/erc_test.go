package tokens

import (
	_ "embed"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

//go:embed erc_mock_balance.json
var mock_balance string

//go:embed erc_mock_details.json
var mock_details string

func TestErcDiscovery(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/balance/"):
			w.Write([]byte(mock_balance))
		case strings.Contains(r.URL.Path, "/details/"):
			w.Write([]byte(mock_details))
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	cmcService := &CMCService{
		logger:     logrus.WithField("module", "cmc_id_service").Logger,
		cachedData: cache.New(10*time.Hour, 1*time.Hour),
	}
	lruCache, err := lru.New[string, models.CoinBase](20000)
	if err != nil {
		t.Errorf("failed to create LRU cache: %v", err)
	}
	oneInchService := &oneInchService{
		logger:     logrus.WithField("module", "oneInch_service").Logger,
		cachedData: lruCache,
	}
	dicoveryService := &ercDiscoveryService{
		logger:         logrus.WithField("module", "oneInch_evm_base_service").Logger,
		baseAddress:    mockServer.URL,
		cmcService:     cmcService,
		oneInchService: oneInchService,
	}
	dicoveryService.oneInchService.oneInchBaseURL = mockServer.URL
	dicoveryService.cmcService.cachedData.Set(cmcService.getCacheKey("Ethereum", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), 2396, cache.DefaultExpiration)
	dicoveryService.oneInchService.cachedData.Add(oneInchService.getCacheKey("Ethereum", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), models.CoinBase{
		Decimals:        18,
		Ticker:          "WETH",
		ContractAddress: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
	})

	dicoveryService.cmcService.cachedData.Set(cmcService.getCacheKey("Ethereum", "0xdac17f958d2ee523a2206206994597c13d831ec7"), 825, cache.DefaultExpiration)
	dicoveryService.oneInchService.cachedData.Add(oneInchService.getCacheKey("Ethereum", "0xdac17f958d2ee523a2206206994597c13d831ec7"), models.CoinBase{
		Decimals:        6,
		Ticker:          "Tether",
		ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
	})

	res, err := dicoveryService.Discover("0x14F6Ed6CBb27b607b0E2A48551A988F1a19c89B6", common.Ethereum)
	if err != nil {
		t.Errorf("oneInchEVMBase failed: %v", err)
	}

	expectedTokens := 2 // Based on the mock response data
	if len(res) != expectedTokens {
		t.Errorf("Expected %d tokens, got %d", expectedTokens, len(res))
	}

	// Test for specific token presence (WETH as an example)
	found := false
	for _, token := range res {
		if token.Ticker == "WETH" {
			found = true
			if token.CMCId != 2396 {
				t.Errorf("Expected WETH CMC ID to be 2396, got %d", token.CMCId)
			}
			break
		}
	}
	if !found {
		t.Error("Expected to find WETH token in results")
	}
	found = false
	for _, token := range res {
		if token.Ticker == "Tether" {
			found = true
			if token.CMCId != 825 {
				t.Errorf("Expected Tether CMC ID to be 825, got %d", token.CMCId)
			}
			break
		}
	}
	if !found {
		t.Error("Expected to find Tether token in results")
	}
}

func Test_ercDiscoveryService_Search(t *testing.T) {
	// Setup test server
	mockServer := newMockServer()
	defer mockServer.Close()

	// Initialize discovery service
	discovery := setupDiscoveryService(mockServer.URL)

	type CoinInfo struct {
		CoinBase models.CoinBase
		CMCId    int
	}
	// Test cases
	testCases := []struct {
		name     string
		input    models.CoinBase
		expected CoinInfo
	}{
		{
			name: "WETH token search",
			input: models.CoinBase{
				Ticker:          "WETH",
				ContractAddress: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
				Decimals:        18,
				Chain:           common.Ethereum,
			},
			expected: CoinInfo{
				CoinBase: models.CoinBase{
					Ticker:          "WETH",
					ContractAddress: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
					Decimals:        18,
					Chain:           common.Ethereum,
				},
				CMCId: 2396,
			},
		},
		{
			name: "Tether token search",
			input: models.CoinBase{
				Ticker:          "Tether",
				ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
				Decimals:        6,
				Chain:           common.Ethereum,
			},
			expected: CoinInfo{
				CoinBase: models.CoinBase{
					Ticker:          "Tether",
					ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
					Decimals:        6,
					Chain:           common.Ethereum,
				},
				CMCId: 825,
			},
		},
	}

	// Run tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := discovery.Search(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Ticker != tc.expected.CoinBase.Ticker {
				t.Errorf("expected ticker %s, got %s", tc.expected.CoinBase.Ticker, result.Ticker)
			}
			if result.CMCId != tc.expected.CMCId {
				t.Errorf("expected CMC ID %d, got %d", tc.expected.CMCId, result.CMCId)
			}
			if result.ContractAddress != tc.expected.CoinBase.ContractAddress {
				t.Errorf("expected contract address %s, got %s", tc.expected.CoinBase.ContractAddress, result.ContractAddress)
			}
			if result.Decimals != tc.expected.CoinBase.Decimals {
				t.Errorf("expected decimals %d, got %d", tc.expected.CoinBase.Decimals, result.Decimals)
			}
			if result.Chain != tc.expected.CoinBase.Chain {
				t.Errorf("expected chain %s, got %s", tc.expected.CoinBase.Chain, result.Chain)
			}
		})
	}
}

func newMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/balance/"):
			w.Write([]byte(mock_balance))
		case strings.Contains(r.URL.Path, "/details/"):
			w.Write([]byte(mock_details))
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}))
}

func setupDiscoveryService(baseURL string) *ercDiscoveryService {
	cmcService := &CMCService{
		logger:     logrus.WithField("module", "cmc_id_service").Logger,
		cachedData: cache.New(10*time.Hour, 1*time.Hour),
	}
	cache, err := lru.New[string, models.CoinBase](20000)
	if err != nil {
		log.Panic("failed to create LRU cache: ", err)
	}

	oneInchService := &oneInchService{
		logger:     logrus.WithField("module", "oneInch_service").Logger,
		cachedData: cache,
	}
	discovery := &ercDiscoveryService{
		logger:         logrus.WithField("module", "oneInch_evm_base_service").Logger,
		baseAddress:    baseURL,
		cmcService:     cmcService,
		oneInchService: oneInchService,
	}

	// Setup cache
	discovery.oneInchService.oneInchBaseURL = baseURL
	setupTestCache(discovery)

	return discovery
}

func setupTestCache(discovery *ercDiscoveryService) {
	tokens := map[string]struct {
		address  string
		cmcId    int
		ticker   string
		decimals int
	}{
		"WETH": {
			address:  "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			cmcId:    2396,
			ticker:   "WETH",
			decimals: 18,
		},
		"Tether": {
			address:  "0xdac17f958d2ee523a2206206994597c13d831ec7",
			cmcId:    825,
			ticker:   "Tether",
			decimals: 6,
		},
	}

	for _, token := range tokens {
		discovery.cmcService.cachedData.Set(
			discovery.cmcService.getCacheKey("Ethereum", token.address),
			token.cmcId,
			cache.DefaultExpiration,
		)
		discovery.oneInchService.cachedData.Add(
			discovery.oneInchService.getCacheKey("Ethereum", token.address),
			models.CoinBase{
				Decimals:        token.decimals,
				Ticker:          token.ticker,
				ContractAddress: token.address,
			},
			
		)
	}
}
