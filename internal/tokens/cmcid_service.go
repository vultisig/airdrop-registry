package tokens

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

var cmcChainMap map[common.Chain]string = map[common.Chain]string{
	common.Bitcoin:      "Bitcoin",
	common.Ethereum:     "Ethereum",
	common.THORChain:    "THORChain",
	common.Solana:       "Solana",
	common.Avalanche:    "Avalanche",
	common.BscChain:     "BNB",
	common.BitcoinCash:  "Bitcoin Cash",
	common.Litecoin:     "Litecoin Cash",
	common.Dogecoin:     "Dogecoin",
	common.GaiaChain:    "GaiaChain",
	common.Kujira:       "Kujira",
	common.Dash:         "Dash",
	common.MayaChain:    "MayaChain",
	common.Arbitrum:     "Arbitrum",
	common.Base:         "Base",
	common.Optimism:     "Optimism",
	common.Polygon:      "POL (prev. MATIC)",
	common.Blast:        "Blast",
	common.CronosChain:  "CronosChain",
	common.Sui:          "Sui",
	common.Polkadot:     "Polkadot",
	common.Zksync:       "Zksync",
	common.Dydx:         "Dydx",
	common.Ton:          "Toncoin",
	common.Terra:        "Terra",
	common.TerraClassic: "TerraClassic",
	common.XRP:          "XRP",
	common.Osmosis:      "Osmosis",
	common.Noble:        "NOBLEBLOCKS",
	common.Tron:         "TRON",
}

type contractCoin struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type platform struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	TokenAddress string       `json:"token_address"`
	Coin         contractCoin `json:"coin"`
}
type contract struct {
	ContractAddress string    `json:"contract_address"`
	Platform        *platform `json:"platform"`
}
type contractToken struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	ContractAddresses []contract `json:"contract_address"`
}
type cmcAsset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type mainData struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Platform *platform `json:"platform"`
}
type contractModel struct {
	Data map[string]contractToken `json:"data"`
}
type mainModel struct {
	Data []mainData `json:"data"`
}

type CMCIDService struct {
	logger     *logrus.Logger
	CMCBaseURL string
	cachedData *cache.Cache
}

var cmcNativeCoins map[string]int

func NewCMCIDService() *CMCIDService {
	cmcNativeCoins = make(map[string]int)
	cmcBaseURL := "https://api.vultisig.com/cmc/v1/cryptocurrency"
	cachedData := cache.New(10*time.Hour, 1*time.Hour)
	var cmcMainModel mainModel
	start, limit := 1, 5000
	for {
		url := fmt.Sprintf("%s/map?sort=cmc_rank&limit=%d&start=%d", cmcBaseURL, limit, start)
		resp, err := http.Get(url)
		if err != nil {
			logrus.Errorf("error fetching cmc id from %s: %e", url, err)
			return nil
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			logrus.Errorf("failed to get data from %s, status code: %d", url, resp.StatusCode)
			return nil
		}
		if err := json.NewDecoder(resp.Body).Decode(&cmcMainModel); err != nil {
			logrus.Errorf("error decoding cmc id from %s: %e", url, err)
			return nil
		}
		for _, v := range cmcMainModel.Data {
			key := v.Name
			if v.Platform != nil {
				key = v.Platform.Name + v.Platform.TokenAddress
			}
			cmcNativeCoins[key] = v.ID
		}
		if len(cmcMainModel.Data) < limit {
			break
		}
		start += limit
	}
	return &CMCIDService{
		logger:     logrus.WithField("module", "cmc_id_service").Logger,
		CMCBaseURL: cmcBaseURL,
		cachedData: cachedData,
	}
}

func (c *CMCIDService) GetCMCID(chain common.Chain, coin models.Coin) (int, error) {
	if coin.ContractAddress == "" {
		if cmcID, ok := cmcNativeCoins[cmcChainMap[chain]]; ok {
			return cmcID, nil
		}
	}
	return c.GetCMCIDByContract(cmcChainMap[chain], coin.ContractAddress)
}

func (c *CMCIDService) GetCMCIDByContract(chain, contract string) (int, error) {
	key := chain + contract
	if cachedData, found := c.cachedData.Get(key); found {
		if cmcID, ok := cachedData.(int); ok {
			return cmcID, nil
		}
	}
	url := fmt.Sprintf("%s/info?address=%s&skip_invalid=true&aux=status", c.CMCBaseURL, contract)
	resp, err := http.Get(url)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("failed to get data from %s, status code: %d", url, resp.StatusCode)
	}
	var cmcContractModel contractModel
	if err := json.NewDecoder(resp.Body).Decode(&cmcContractModel); err != nil {
		return -1, err
	}
	for _, v := range cmcContractModel.Data {
		for _, ca := range v.ContractAddresses {
			if ca.ContractAddress == contract && ca.Platform.Coin.Name == chain {
				key := ca.Platform.Coin.Name + contract
				c.cachedData.Set(key, v.ID, cache.DefaultExpiration)
				return v.ID, nil
			}
		}
	}
	return -1, fmt.Errorf("failed to get cmc id for contract: %s", contract)
}
