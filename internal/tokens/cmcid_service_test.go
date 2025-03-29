package tokens

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func TestCMCIDService_GetCMCID(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := mainModel{
			Data: []mainData{
				{
					ID:   1027,
					Name: "Ethereum",
				},
				{
					ID:   228261,
					Name: "Osmosis",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()
	cachedData := cache.New(10*time.Hour, 1*time.Hour)
	cmcService := &CMCService{
		baseURL:    mockServer.URL,
		cachedData: cachedData,
		nativeCoinIds: map[string]int{
			"3DPass":    28794,
			"42-coin":   93,
			"AB":        3871,
			"ABBC Coin": 3437,
			"Osmosisibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877": 22861,
			"Rubix":    17972,
			"Ethereum": 1027,
			"Ethereum0xdac17f958d2ee523a2206206994597c13d831ec7": 825,
			"BNB0xA697e272a73744b343528C3Bc4702F2565b2F422":      23095,
		},
	}
	cmcService.baseURL = mockServer.URL
	type v int
	cmcService.cachedData.Set(cmcService.getCacheKey(cmcChainMap[common.Osmosis], "ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877"), 228261, cache.DefaultExpiration)
	type cmc struct {
		chain         common.Chain
		asset         models.Coin
		expectedCMCID int
	}
	cmcVals := []cmc{
		cmc{
			chain:         common.Ethereum,
			expectedCMCID: 1027,
		},
		cmc{
			chain: common.Osmosis,
			asset: models.Coin{
				ContractAddress: "ibc/D79E7D83AB399BFFF93433E54FAA480C191248FC556924A2A8351AE2638B3877",
			},
			expectedCMCID: 228261,
		},
	}

	for _, cmc := range cmcVals {
		cmcid, err := cmcService.GetCMCID(cmc.chain, cmc.asset)
		assert.NoError(t, err)
		assert.Equal(t, cmc.expectedCMCID, cmcid)
	}

}
