package liquidity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func TestGetSaverPosition(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := map[string][]interface{}{
			"pools": {
				map[string]interface{}{
					"assetAdded":     "4",
					"assetAddress":   "0x1d204941ca5ff1143caca57d71ead1179ba1dd3a",
					"assetDeposit":   "547855400",
					"assetRedeem":    "1200000000",
					"assetWithdrawn": "0",
					"dateFirstAdded": "1726804546",
					"dateLastAdded":  "1726804546",
					"pool":           "AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E",
					"saverUnits":     "499835103",
				},
				map[string]interface{}{
					"assetAdded":     "3",
					"assetAddress":   "0x3d204941ca5ff1143caca57d71ead1179ba1dd3a",
					"assetDeposit":   "527592",
					"assetRedeem":    "1300000000",
					"assetWithdrawn": "0",
					"dateFirstAdded": "1712193841",
					"dateLastAdded":  "1712193841",
					"pool":           "BSC.BNB",
					"saverUnits":     "477850",
				},
				map[string]interface{}{
					"assetAdded":     "2",
					"assetAddress":   "0x3d204941ca5ff1143caca57d71ead1179ba1dd3a",
					"assetDeposit":   "219100",
					"assetRedeem":    "1400000000",
					"assetWithdrawn": "0",
					"dateFirstAdded": "1726459134",
					"dateLastAdded":  "1726459134",
					"pool":           "ETH.ETH",
					"saverUnits":     "200809",
				},
				map[string]interface{}{
					"assetAdded":     "1",
					"assetAddress":   "0x3d204941ca5ff1143caca57d71ead1179ba1dd3a",
					"assetDeposit":   "165316700",
					"assetRedeem":    "1500000000",
					"assetWithdrawn": "500030745",
					"dateFirstAdded": "1701211207",
					"dateLastAdded":  "1725317140",
					"pool":           "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7",
					"saverUnits":     "316998630",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	poolCache := cache.New(5*time.Minute, 10*time.Minute)
	// Create a LiquidityPositionResolver instance
	saverPositionResolver := &SaverPositionResolver{
		midgardBaseURL: mockServer.URL,
		poolCache:      poolCache,
	}
	poolCache.Add("AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E", poolResp{
		Pool:          "AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E",
		AssetPriceUsd: 1}, cache.DefaultExpiration)

	poolCache.Add("BSC.BNB", poolResp{
		Pool:          "BSC.BNB",
		AssetPriceUsd: 300}, cache.DefaultExpiration)

	poolCache.Add("ETH.ETH", poolResp{
		Pool:          "ETH.ETH",
		AssetPriceUsd: 2000}, cache.DefaultExpiration)
	poolCache.Add("ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", poolResp{
		Pool:          "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7",
		AssetPriceUsd: 1}, cache.DefaultExpiration)
	addrs := []string{"0x3d204941ca5ff1143caca57d71ead1179ba1dd3a", "0x1d204941ca5ff1143caca57d71ead1179ba1dd3a"}
	position, err := saverPositionResolver.GetSaverPosition(strings.Join(addrs, ","))
	assert.NoErrorf(t, err, "Failed to get saver position: %v", err)
	assert.Equal(t, float64(31927), position)
}

func TestFetchPools(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := []map[string]interface{}{
			{
				"pool":                      "AVAX.AVAX",
				"assetDepth":                "15977431199949",
				"runeOrCacaoDepth":          "72955181892989",
				"assetPriceUsd":             "26.16692121097031",
				"status":                    "available",
				"runeOrCacaoPrice":          "5.73064410936575104904",
				"provider":                  "THORCHAIN",
				"runeOrCacaoLiquidityInUsd": "418080183362764.31556206573977118056",
				"assetLiquidityInUsd":       "418080183362764",
				"totalLiquidityInUsd":       "836160366725529",
				"apr":                       "0.17383818287522756",
				"history":                   []interface{}{},
			},
			{
				"pool":                      "AVAX.SOL-0XFE6B19286885A4F7F55ADAD09C3CD1F906D2478F",
				"assetDepth":                "99451394308",
				"runeOrCacaoDepth":          "3082872535319",
				"assetPriceUsd":             "177.1",
				"status":                    "available",
				"runeOrCacaoPrice":          "5.73064410936575050558",
				"provider":                  "THORCHAIN",
				"runeOrCacaoLiquidityInUsd": "17666845334451.28397420262065658002",
				"assetLiquidityInUsd":       "17666845334451",
				"totalLiquidityInUsd":       "35333690668903",
				"apr":                       "0.3506148600985919",
				"history":                   []interface{}{},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	poolCache := cache.New(5*time.Minute, 10*time.Minute)
	// Create a LiquidityPositionResolver instance
	saverPositionResolver := &SaverPositionResolver{
		thorwalletBaseURL: mockServer.URL,
		poolCache:         poolCache,
	}

	pool, err := saverPositionResolver.getpool("AVAX.SOL-0XFE6B19286885A4F7F55ADAD09C3CD1F906D2478F")
	assert.NoErrorf(t, err, "Failed to get pools: %v", err)
	assert.Equal(t, float64(177.1), pool.AssetPriceUsd)

	assert.Equal(t, 2, saverPositionResolver.poolCache.ItemCount())

	_, found := saverPositionResolver.poolCache.Get("AVAX.AVAX")
	assert.Equal(t, true, found)

	_, found = saverPositionResolver.poolCache.Get("BTC.BTC")
	assert.Equal(t, false, found)
}
