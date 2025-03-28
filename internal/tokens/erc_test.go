package tokens

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

//go:embed erc_mock_balance.json
var mock_balance string

//go:embed erc_mock_details.json
var mock_details string

//go:embed erc_mock_cmcid_response.json
var mock_cmcid_response string

func TestErcDiscovery(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/balance/"):
			w.Write([]byte(mock_balance))
		case strings.Contains(r.URL.Path, "/details/"):
			w.Write([]byte(mock_details))
		case strings.Contains(r.URL.Path, "info"):
			w.Write([]byte(mock_cmcid_response))
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	oneInchevmbaseservice := &ercDiscoveryService{
		logger:         logrus.WithField("module", "oneInch_evm_base_service").Logger,
		baseAddress:    mockServer.URL,
		cmcIDService:   NewCMCIDService(),
		oneinchService: NewOneinchService(),
	}
	oneInchevmbaseservice.oneinchService.oneinchBaseURL = mockServer.URL
	oneInchevmbaseservice.cmcIDService.cachedData.Set("Ethereum0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", 2396, cache.DefaultExpiration)
	oneInchevmbaseservice.oneinchService.cachedData.Set("Ethereum0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", models.CoinBase{
		Decimals:        18,
		Ticker:          "WETH",
		ContractAddress: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
	}, cache.DefaultExpiration)

	res, err := oneInchevmbaseservice.discover("0x14F6Ed6CBb27b607b0E2A48551A988F1a19c89B6", common.Ethereum)
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
