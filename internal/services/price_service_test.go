package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
