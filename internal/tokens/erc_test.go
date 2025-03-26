package tokens

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func TestOneInchEVMBaseService(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		balance := map[string]any{
			"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2": "1064865514300467",
			"0xc28e931814725bbeb9e670676fabbcb694fe7df2": "0",
			"0xdac17f958d2ee523a2206206994597c13d831ec7": "2088424",
			"0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee": "4566634655444195",
		}
		details := map[string]any{
			"assets": map[string]any{
				"name":     "WETH",
				"symbol":   "WETH",
				"type":     "ERC20",
				"decimals": 18,
				"website":  "https://weth.io/",
				"explorer": "https://etherscan.io/token/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				"status":   "active",
				"id":       "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			},
		}

		cmcIdResponse := map[string]any{
			"data": []map[string]any{
				{
					"id":                    2396,
					"rank":                  10446,
					"name":                  "WETH",
					"symbol":                "WETH",
					"slug":                  "weth",
					"is_active":             1,
					"status":                1,
					"first_historical_data": "2018-01-14T19:05:00.000Z",
					"last_historical_data":  "2025-03-15T14:55:00.000Z",
					"platform": map[string]any{
						"id":            1,
						"name":          "Ethereum",
						"symbol":        "ETH",
						"slug":          "ethereum",
						"token_address": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
					},
				},
				{
					"id":                    31312,
					"rank":                  9245,
					"name":                  "Popo",
					"symbol":                "POPO",
					"slug":                  "pepe-popo",
					"is_active":             1,
					"status":                1,
					"first_historical_data": "2024-05-20T02:10:00.000Z",
					"last_historical_data":  "2025-03-15T15:00:00.000Z",
					"platform": map[string]any{
						"id":            1,
						"name":          "Ethereum",
						"symbol":        "ETH",
						"slug":          "ethereum",
						"token_address": "0x195be8ee12aa1591902c4232b5b25017a9cbbdea",
					},
				},
			},
		}
		switch {
		case strings.Contains(r.URL.Path, "/details/"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(details)
		case strings.Contains(r.URL.Path, "/balance/"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(balance)
		case strings.Contains(r.URL.Path, "info"):
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(cmcIdResponse)
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
