package tokens

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"time"

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
	oneInchService := &oneInchService{
		logger:     logrus.WithField("module", "oneInch_service").Logger,
		cachedData: cache.New(10*time.Hour, 1*time.Hour),
	}
	oneInchevmbaseservice := &ercDiscoveryService{
		logger:         logrus.WithField("module", "oneInch_evm_base_service").Logger,
		baseAddress:    mockServer.URL,
		cmcService:     cmcService,
		oneInchService: oneInchService,
	}
	oneInchevmbaseservice.oneInchService.oneInchBaseURL = mockServer.URL
	oneInchevmbaseservice.cmcService.cachedData.Set(cmcService.getCacheKey("Ethereum", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), 2396, cache.DefaultExpiration)
	oneInchevmbaseservice.oneInchService.cachedData.Set(oneInchService.getCacheKey("Ethereum", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), models.CoinBase{
		Decimals:        18,
		Ticker:          "WETH",
		ContractAddress: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
	}, cache.DefaultExpiration)

	oneInchevmbaseservice.cmcService.cachedData.Set(cmcService.getCacheKey("Ethereum", "0xdac17f958d2ee523a2206206994597c13d831ec7"), 825, cache.DefaultExpiration)
	oneInchevmbaseservice.oneInchService.cachedData.Set(oneInchService.getCacheKey("Ethereum", "0xdac17f958d2ee523a2206206994597c13d831ec7"), models.CoinBase{
		Decimals:        6,
		Ticker:          "Tether",
		ContractAddress: "0xdac17f958d2ee523a2206206994597c13d831ec7",
	}, cache.DefaultExpiration)

	res, err := oneInchevmbaseservice.Discover("0x14F6Ed6CBb27b607b0E2A48551A988F1a19c89B6", common.Ethereum)
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
}
