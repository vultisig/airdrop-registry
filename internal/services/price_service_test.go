package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetCMCMap(t *testing.T) {
	pr, err := NewPriceResolver()
	assert.Nil(t, err)
	assert.NotNil(t, pr)
}

func TestGetLifiPrice(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := map[string]interface{}{
			"address":  "0x815C23eCA83261b6Ec689b60Cc4a58b54BC24D8D",
			"chainId":  1,
			"symbol":   "vTHOR",
			"decimals": 18,
			"name":     "vTHOR",
			"coinKey":  "vTHOR",
			"priceUSD": "1.5",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a PriceResolver instance
	priceResolver := &PriceResolver{
		lifiBaseAddress: mockServer.URL,
	}

	price, err := priceResolver.GetLiFiPrice("eth", "0x815C23eCA83261b6Ec689b60Cc4a58b54BC24D8D")
	if err != nil {
		t.Fatalf("Failed to get VThor price: %v", err)
	}

	if price != 1.5 {
		t.Errorf("Expected price 100.5 but got %v", price)
	}
}

func TestCacaoPrice(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := map[string]map[string]interface{}{"cacao": {"usd": 0.5}}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a PriceResolver instance
	priceResolver := &PriceResolver{
		logger:               logrus.WithField("module", "price_resolver").Logger,
		priceCache:           *cache.New(4*time.Minute, 5*time.Minute),
		coingeckoBaseAddress: mockServer.URL,
	}

	price, err := priceResolver.GetCoinGeckoPrice("cacao", "usd")
	assert.NoErrorf(t, err, "Failed to get CACAO price: %v", err)
	assert.Equal(t, float64(0.5), price)

	price, err = priceResolver.GetCoinGeckoPrice("CACAO", "USD")
	assert.EqualError(t, err, "price not found in response")
	assert.Equal(t, float64(0), price)
}
