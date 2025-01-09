package balance

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFetchThorchainBalanceOfAddress(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := map[string]interface{}{
			"balances": []map[string]interface{}{
				{
					"denom":  "rune",
					"amount": "2500000000",
				},
			},
			"pagination": map[string]interface{}{
				"next_key": nil,
				"total":    "1",
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a LiquidityPositionResolver instance
	balanceResolver := &BalanceResolver{
		logger:                 logrus.WithField("module", "balance_resolver_test").Logger,
		thornodeBaseAddress:    mockServer.URL,
		thorchainRuneProviders: &sync.Map{},
		thorchainBondProviders: &sync.Map{},
	}
	balanceResolver.thorchainRuneProviders.Store("thor2rjxghep6g3j3z0k3jwz3wzrj3z0k3jwz3wzrj", int64(1200000000))
	balanceResolver.thorchainBondProviders.Store("thor2rjxghep6g3j3z0k3jwz3wzrj3z0k3jwz3wzrj", "2000000000")
	balance, err := balanceResolver.FetchThorchainBalanceOfAddress("thor2rjxghep6g3j3z0k3jwz3wzrj3z0k3jwz3wzrj")
	assert.NoErrorf(t, err, "Failed to get thorchain rune providers: %v", err)
	assert.Equal(t, float64(57), balance)

	balance, err = balanceResolver.FetchThorchainBalanceOfAddress("thor2")
	assert.NoErrorf(t, err, "Failed to get thorchain rune providers: %v", err)
	assert.Equal(t, float64(25), balance)
}

func TestGetTHORChainRuneProviders(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := []map[string]interface{}{
			{
				"rune_address":         "thor2rjxghep6g3j3z0k3jwz3wzrj3z0k3jwz3wzrj",
				"units":                "29856738754",
				"value":                "28892079482",
				"pnl":                  "-1107920518",
				"deposit_amount":       "30000000000",
				"withdraw_amount":      "0",
				"last_deposit_height":  18123258,
				"last_withdraw_height": 0,
			},
			{
				"rune_address":         "thor1cfzgzg02cp7yjrkagzdrdp7dqh0xlsdhawwjc",
				"units":                "24511378597",
				"value":                "23719425771",
				"pnl":                  "-1019445922",
				"deposit_amount":       "109894000000",
				"withdraw_amount":      "85155128307",
				"last_deposit_height":  17571421,
				"last_withdraw_height": 17571407,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a LiquidityPositionResolver instance
	balanceResolver := &BalanceResolver{
		logger:                 logrus.WithField("module", "balance_resolver_test").Logger,
		thornodeBaseAddress:    mockServer.URL,
		thorchainRuneProviders: &sync.Map{},
	}
	err := balanceResolver.GetTHORChainRuneProviders()
	assert.NoErrorf(t, err, "Failed to get thorchain rune providers: %v", err)

	value, ok := balanceResolver.thorchainRuneProviders.Load("thor1cfzgzg02cp7yjrkagzdrdp7dqh0xlsdhawwjc")
	assert.True(t, ok)
	assert.Equal(t, int64(23719425771), value)

}

func TestFetchTerraBalanceOfAddress(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := map[string]interface{}{
			"balances": []map[string]interface{}{
				{
					"denom":  "uluna",
					"amount": "2500000000",
				},
			},
			"pagination": map[string]interface{}{
				"next_key": nil,
				"total":    "1",
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a LiquidityPositionResolver instance
	balanceResolver := &BalanceResolver{
		logger: logrus.WithField("module", "balance_resolver_test").Logger,
	}
	balance, err := balanceResolver.fetchSpecificCosmosBalance(mockServer.URL+"/cosmos/bank/v1beta1/spendable_balances/"+"terra1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3nln0mh", "uluna", 6)
	assert.NoErrorf(t, err, "Failed to get thorchain rune providers: %v", err)
	assert.Equal(t, float64(2500), balance)
}
