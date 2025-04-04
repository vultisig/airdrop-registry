package liquidity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetLiquidityPosition(t *testing.T) {
	// Set up mock server to handle both Thor and Maya endpoints
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tcMemberresponse := map[string]any{
			"pools": []map[string]any{
				{
					"pool":         "AVAX.AVAX",
					"assetDeposit": "1182138963",
					"runeDeposit":  "13085493023",
				},
			},
		}
		mayaMemberresponse := map[string]any{
			"pools": []map[string]any{
				{
					"pool":         "THOR.RUNE",
					"assetDeposit": "0",
					"cacaoDeposit": "0",
				},
				{
					"pool":         "THOR.RUNE",
					"assetDeposit": "0",
					"cacaoDeposit": "30000000000",
				},
			},
		}

		switch {
		case strings.Contains(r.URL.Path, "maya/v2/member"):
			json.NewEncoder(w).Encode(mayaMemberresponse)
		case strings.Contains(r.URL.Path, "thor/v2/member"):
			json.NewEncoder(w).Encode(tcMemberresponse)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer mockServer.Close()

	// Initialize resolvers with test data
	tcusdPools := []string{
		"ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7",
		"AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E",
		"ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
	}

	tclpResolver := &MidgardLPResolver{
		baseAddress: mockServer.URL + "/thor",
		usdPools:    tcusdPools,
		poolCache:   cache.New(3*time.Hour, 6*time.Hour),
		logger:      logrus.WithField("module", "tc_position_resolver").Logger,
	}
	mayalpResolver := &MidgardLPResolver{
		baseAddress: mockServer.URL + "/maya",
		usdPools: []string{
			"ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7",
			"ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
		},
		poolCache: cache.New(3*time.Hour, 6*time.Hour),
		logger:    logrus.WithField("module", "maya_position_resolver").Logger,
	}

	// Populate cache with test pool data
	tcPoolsResponse := []pool{
		{Asset: "AVAX.AVAX", AssetPrice: 15.986349195651876, AssetPriceUSD: 17.76209190099236},
		{Asset: "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", AssetPrice: 0.9000262629402749, AssetPriceUSD: 1},
		{Asset: "AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E", AssetPrice: 0.87508425109807, AssetPriceUSD: 0.9999999999999999},
		{Asset: "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48", AssetPrice: 0.8752408201484536, AssetPriceUSD: 1.0001789188300294},
	}
	mayaPoolsResponse := []pool{
		{Asset: "THOR.RUNE", AssetPrice: 8.528852444728814, AssetPriceUSD: 1.168208562804459},
		{Asset: "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", AssetPrice: 7.48114860381894, AssetPriceUSD: 1},
		{Asset: "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48", AssetPrice: 7.3442966936902785, AssetPriceUSD: 1.0004148785977214},
	}

	for _, p := range tcPoolsResponse {
		tclpResolver.poolCache.Set(p.Asset, p, cache.DefaultExpiration)
	}
	for _, p := range mayaPoolsResponse {
		mayalpResolver.poolCache.Set(p.Asset, p, cache.DefaultExpiration)
	}

	// Test GetLiquidityPosition
	liquidityPositionResolver := NewLiquidtyPositionResolver()
	liquidityPositionResolver.tclpResolver = tclpResolver
	liquidityPositionResolver.mayalpResolver = mayalpResolver

	addrs := []string{
		"maya1zmg5l34d6sf0xk7rwwjnz45vjs5pa7cwduwac8",
		"thor1005rk5k9uuew3u5y489yd8tgjyrsykknnat8z0",
	}

	lp, err := liquidityPositionResolver.GetLiquidityPosition(strings.Join(addrs, ","))
	assert.NoError(t, err)
	assert.Equal(t, float64(39860826147.002235), lp)
}

func TestGetMemberPosition(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "tc/v2/member/"):
			tcMemberresponse := map[string]any{
				"pools": []map[string]any{
					{
						"pool":         "AVAX.AVAX",
						"assetDeposit": "1182138963",
						"runeDeposit":  "13085493023",
					},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tcMemberresponse)
		case strings.Contains(r.URL.Path, "maya/v2/member/"):
			mayaMemberresponse := map[string]any{
				"pools": []map[string]any{
					{
						"pool":         "THOR.RUNE",
						"assetDeposit": "0",
						"cacaoDeposit": "0",
					},
					{
						"pool":         "THOR.RUNE",
						"assetDeposit": "0",
						"cacaoDeposit": "30000000000",
					},
				},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(mayaMemberresponse)
		case strings.Contains(r.URL.Path, "/v2/member/"):
			w.WriteHeader(http.StatusNotFound)
		case strings.Contains(r.URL.Path, "/v2/member/invalidAddress"):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	tcResolver := &MidgardLPResolver{
		baseAddress: mockServer.URL + "/tc/v2/member",
		logger:      logrus.WithField("module", "test").Logger,
	}

	// Test valid address
	members, err := tcResolver.GetMemberPosition("maya1zmg5l34d6sf0xk7rwwjnz45vjs5pa7cwduwac8,thor1005rk5k9uuew3u5y489yd8tgjyrsykknnat8z0")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(members.Pools))
	assert.Equal(t, float64(1182138963), members.Pools[0].AssetDeposit)
	assert.Equal(t, "AVAX.AVAX", members.Pools[0].Pool)
	assert.Equal(t, int64(13085493023), members.Pools[0].RuneDeposit)

	// Test empty address returns empty members
	members, err = tcResolver.GetMemberPosition("")
	assert.NoError(t, err)

	// Test server error
	tcResolver.baseAddress = "http://invalid-url"
	members, err = tcResolver.GetMemberPosition("maya1zmg5l34d6sf0xk7rwwjnz45vjs5pa7cwduwac8,thor1005rk5k9uuew3u5y489yd8tgjyrsykknnat8z0")
	assert.Error(t, err)
	assert.Empty(t, members.Pools)

	mayaResolver := &MidgardLPResolver{
		baseAddress: mockServer.URL + "/maya/v2/member",
		logger:      logrus.WithField("module", "test").Logger,
	}
	members, err = mayaResolver.GetMemberPosition("maya1zmg5l34d6sf0xk7rwwjnz45vjs5pa7cwduwac8,thor1005rk5k9uuew3u5y489yd8tgjyrsykknnat8z0")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(members.Pools))
	assert.Equal(t, float64(0), members.Pools[0].AssetDeposit)
	assert.Equal(t, float64(0), members.Pools[1].AssetDeposit)
	assert.Equal(t, "THOR.RUNE", members.Pools[0].Pool)
	assert.Equal(t, "THOR.RUNE", members.Pools[1].Pool)
	assert.Equal(t, int64(0), members.Pools[0].CacaoDeposit)
	assert.Equal(t, int64(30000000000), members.Pools[1].CacaoDeposit)

	// Test empty address returns empty members
	members, err = mayaResolver.GetMemberPosition("")
	assert.NoError(t, err)

	// Test server error
	mayaResolver.baseAddress = "http://invalid-url"
	members, err = mayaResolver.GetMemberPosition("validAddress")
	assert.Error(t, err)
	assert.Empty(t, members.Pools)

}

func TestRefreshCache(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/pools" {
			tcPoolsResponse := []pool{
				{Asset: "AVAX.AVAX", AssetPrice: 15.986349195651876, AssetPriceUSD: 17.76209190099236},
				{Asset: "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", AssetPrice: 0.9000262629402749, AssetPriceUSD: 1},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tcPoolsResponse)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	resolver := &MidgardLPResolver{
		baseAddress: mockServer.URL,
		poolCache:   cache.New(3*time.Hour, 6*time.Hour),
		logger:      logrus.WithField("module", "test").Logger,
	}

	// Test successful cache refresh
	err := resolver.refreshCache()
	assert.NoError(t, err)

	// Verify cache was populated
	avaxPool, found := resolver.poolCache.Get("AVAX.AVAX")
	assert.True(t, found)
	assert.Equal(t, "AVAX.AVAX", avaxPool.(pool).Asset)
	assert.Equal(t, 17.76209190099236, avaxPool.(pool).AssetPriceUSD)

	ethPool, found := resolver.poolCache.Get("ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7")
	assert.True(t, found)
	assert.Equal(t, "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", ethPool.(pool).Asset)
	assert.Equal(t, float64(1), ethPool.(pool).AssetPriceUSD)

	// Test with bad server URL
	resolver.baseAddress = "http://invalid-url"
	err = resolver.refreshCache()
	assert.Error(t, err)
}

func TestGetAssetPrice(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/pools" {
			pools := []pool{
				{Asset: "AVAX.AVAX", AssetPrice: 15.986349195651876, AssetPriceUSD: 17.76209190099236},
				{Asset: "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", AssetPrice: 0.9000262629402749, AssetPriceUSD: 1},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(pools)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	resolver := &MidgardLPResolver{
		baseAddress: mockServer.URL,
		poolCache:   cache.New(3*time.Hour, 6*time.Hour),
		logger:      logrus.WithField("module", "test").Logger,
	}

	// Test getting price from cache
	resolver.poolCache.Set("AVAX.AVAX", pool{
		Asset:         "AVAX.AVAX",
		AssetPrice:    15.986349195651876,
		AssetPriceUSD: 17.76209190099236,
	}, cache.DefaultExpiration)

	price, err := resolver.getAssetPrice("AVAX.AVAX")
	assert.NoError(t, err)
	assert.Equal(t, float64(17.76209190099236), price)

	// Test getting price after refresh
	price, err = resolver.getAssetPrice("ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7")
	assert.NoError(t, err)
	assert.Equal(t, float64(1), price)

	// Test with non-existent asset
	price, err = resolver.getAssetPrice("NONEXISTENT")
	assert.NoError(t, err)
	assert.Equal(t, float64(0), price)

	// Test with bad server URL
	resolver.baseAddress = "http://invalid-url"
	_, err = resolver.getAssetPrice("BTC.BTC")
	assert.Error(t, err)
}

func TestGetNativeTokenPrice(t *testing.T) {
	// Initialize resolver with test data
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v2/pools" {
			pools := []pool{
				{Asset: "AVAX.AVAX", AssetPrice: 15.986349195651876, AssetPriceUSD: 17.76209190099236},
				{Asset: "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", AssetPrice: 0.9000262629402749, AssetPriceUSD: 1},
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(pools)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	resolver := &MidgardLPResolver{
		baseAddress: mockServer.URL,
		usdPools:    []string{},
		poolCache:   cache.New(3*time.Hour, 6*time.Hour),
		logger:      logrus.WithField("module", "test").Logger,
	}

	// Test 1: Empty cache wont work for 0 result, due to refreshCache()
	price, err := resolver.getNativeTokenPrice()
	assert.NoError(t, err)
	assert.Equal(t, float64(0), price)

	// Test 2: Single USD pool in cache
	resolver.poolCache.Set("ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", pool{
		Asset:         "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7",
		AssetPrice:    2.0,
		AssetPriceUSD: 1.0,
	}, cache.DefaultExpiration)

	price, err = resolver.getNativeTokenPrice()
	assert.NoError(t, err)
	assert.Equal(t, float64(0.5), price) // 1.0/2.0 = 0.5

	// Test 3: Multiple USD pools in cache
	resolver.poolCache.Set("ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48", pool{
		Asset:         "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48",
		AssetPrice:    4.0,
		AssetPriceUSD: 1.0,
	}, cache.DefaultExpiration)

	price, err = resolver.getNativeTokenPrice()
	assert.NoError(t, err)
	assert.Equal(t, float64(0.375), price) // average of (1.0/2.0 + 1.0/4.0)/2 = 0.375

	// Test 4: Zero price values
	resolver.poolCache.Set("ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7", pool{
		Asset:         "ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7",
		AssetPrice:    0,
		AssetPriceUSD: 1.0,
	}, cache.DefaultExpiration)

	price, err = resolver.getNativeTokenPrice()
	assert.NoError(t, err)
	assert.Equal(t, float64(0.25), price) // Only valid pool used: 1.0/4.0 = 0.25
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
	liquidityPositionResolver.SetTGTPrice(0.5)

	tgtlp, err := liquidityPositionResolver.GetTGTStakePosition("0x143A044e411222F36a0f1E35847eCf2400A0d3Df")
	if err != nil {
		t.Fatalf("Failed to get liquidity position: %e", err)
	}
	assert.Equal(t, 170*0.5+1.5, tgtlp)
}
