package liquidity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLiquidityPosition(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := map[string][]interface{}{
			"maya": {
				map[string]interface{}{
					"assetAdded":               "0",
					"assetAddress":             "0x143a044e497624f46a0f1e35847ecf2400a0d3df",
					"assetDeposit":             "0",
					"assetPending":             "0.000767",
					"assetWithdrawn":           "0",
					"cacaoDeposit":             "0",
					"dateFirstAdded":           "0",
					"dateLastAdded":            "0",
					"liquidityUnits":           "0",
					"pool":                     "ARB.ETH",
					"runeAdded":                "0",
					"runeAddress":              "maya18hxxd8tlfpfykg2jpma6j0dgrdv5vmerz7vpyg",
					"runePending":              "0",
					"runeWithdrawn":            "0",
					"assetPriceUsd":            "2519.9153452524106",
					"runeOrCacaoPricePriceUsd": "0.557590620506003",
					"runeOrCacaoAddedUsd":      "1",
					"assetAddedUsd":            "2",
				},
			},
			"thorchain": {
				map[string]interface{}{
					"assetAdded":               "4.611379",
					"assetAddress":             "0x3d504949ca5ff1143caca57d75ead1179ba1dd3a",
					"assetDeposit":             "0",
					"assetPending":             "0",
					"assetWithdrawn":           "592523300",
					"dateFirstAdded":           "1710463968",
					"dateLastAdded":            "1723767813",
					"liquidityUnits":           "968550073",
					"pool":                     "AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E",
					"runeAdded":                "0.81811096",
					"runeAddress":              "thor1ycgzfpcz93qa0xc392xgk75lfee5vdc59hv3r8",
					"runeDeposit":              "0",
					"runePending":              "0",
					"runeWithdrawn":            "112709476",
					"assetPriceUsd":            "1.0007954544738127",
					"runeOrCacaoPricePriceUsd": "5.656201109085557",
					"runeOrCacaoAddedUsd":      "2",
					"assetAddedUsd":            "3.5",
				},
				map[string]interface{}{
					"assetAdded":               "4.77298",
					"assetAddress":             "0x3d504949ca5ff1143caca57d75ead1179ba1dd3a",
					"assetDeposit":             "0",
					"assetPending":             "0",
					"assetWithdrawn":           "396330500",
					"dateFirstAdded":           "1714613024",
					"dateLastAdded":            "1714613024",
					"liquidityUnits":           "501245983",
					"pool":                     "AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E",
					"runeAdded":                "0",
					"runeAddress":              "",
					"runeDeposit":              "0",
					"runePending":              "0",
					"runeWithdrawn":            "0",
					"assetPriceUsd":            "1.0007954544738127",
					"runeOrCacaoPricePriceUsd": "5.656201109085557",
					"runeOrCacaoAddedUsd":      "1.5",
					"assetAddedUsd":            "4",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a LiquidityPositionResolver instance
	liquidityPositionResolver := &LiquidityPositionResolver{
		thorwalletBaseURL: mockServer.URL,
	}

	lp, err := liquidityPositionResolver.GetLiquidityPosition([]string{"thor21cfzgzg02cp7yjrkagzdrdp7dqh0xlsdhawwjc", "0x3d512341ca1ff1142caca57d75ead1179ba1dd3a"})
	if err != nil {
		t.Fatalf("Failed to get liquidity position: %v", err)
	}

	assert.Equal(t, float64(14), lp)
}

func TestGetTGTStakePosition(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := map[string]interface{}{"stakedAmount": "170", "reward": "1.5"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a PriceResolver instance
	liquidityPositionResolver := &LiquidityPositionResolver{
		thorwalletBaseURL: mockServer.URL,
	}

	tgtlp, err := liquidityPositionResolver.GetTGTStakePosition("0x143A044e411222F36a0f1E35847eCf2400A0d3Df")
	if err != nil {
		t.Fatalf("Failed to get liquidity position: %e", err)
	}
	assert.Equal(t, 171.5, tgtlp)
}
