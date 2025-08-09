package tokens

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/utils"
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

type CMCService struct {
	logger        *logrus.Logger
	baseURL       string
	cachedData    *cache.Cache
	nativeCoinIds map[string]int
}

func NewCMCService() (*CMCService, error) {
	cmcService := CMCService{
		logger:        logrus.WithField("module", "cmc_id_service").Logger,
		baseURL:       "https://api.vultisig.com/cmc/v1/cryptocurrency",
		cachedData:    cache.New(10*time.Hour, 1*time.Hour),
		nativeCoinIds: map[string]int{},
	}
	if err := cmcService.initialise(); err != nil {
		return nil, err
	}
	return &cmcService, nil
}
func (c *CMCService) initialise() error {

	start, limit := 1, 5000
	for {
		dataMap, err := c.fetchCMCMap(start, limit)
		if err != nil {
			return err
		}
		for _, v := range dataMap {
			if v.Platform == nil {
				c.nativeCoinIds[v.Name] = v.ID
			} else {
				c.cachedData.Set(c.getCacheKey(v.Platform.Name, v.Platform.TokenAddress), v.ID, cache.DefaultExpiration)
			}
		}
		if len(dataMap) < limit {
			break
		}
		start += limit
	}
	return nil
}

func (c *CMCService) fetchCMCMap(start, limit int) ([]mainData, error) {
	var cmcMainModel mainModel
	url := fmt.Sprintf("%s/map?sort=cmc_rank&limit=%d&start=%d", c.baseURL, limit, start)
	resp, err := http.Get(url)
	if err != nil {
		c.logger.Errorf("error fetching cmc id from %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.logger.Errorf("failed to get data from %s, status code: %d", url, resp.StatusCode)
		return nil, fmt.Errorf("failed to get data from %s, status code: %d", url, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&cmcMainModel); err != nil {
		c.logger.Errorf("error decoding cmc id from %s: %v", url, err)
		return nil, fmt.Errorf("error decoding cmc id from %s: %v", url, err)
	}
	return cmcMainModel.Data, nil
}

func (c *CMCService) GetCMCID(chain common.Chain, coin models.Coin) (int, error) {
	if coin.ContractAddress == "" { // is native coin
		if cmcID, ok := c.nativeCoinIds[cmcChainMap[chain]]; ok {
			return cmcID, nil
		} else {
			return -1, fmt.Errorf("failed to get cmc id for native coin: %s", cmcChainMap[chain])
		}
	}
	return c.GetCMCIDByContract(cmcChainMap[chain], coin.ContractAddress)
}

func (c *CMCService) GetCMCIDByContract(chainName, contract string) (int, error) {
	if cachedData, found := c.cachedData.Get(c.getCacheKey(chainName, contract)); found {
		if cmcID, ok := cachedData.(int); ok {
			return cmcID, nil
		}
	}
	if contract == "" {
		return 0, fmt.Errorf("empty contract address provided")
	}

	//if chain is EVM, then convert to checksum address
	if utils.IsETHAddress(contract) {
		contractEIP55, err := utils.EIP55Checksum(contract)
		if err != nil {
			return -1, fmt.Errorf("failed to convert contract address to EIP55 checksum: %w", err)
		}
		contract = contractEIP55
	}
	time.Sleep(5 * time.Second) // to avoid hitting the API rate limit

	url := fmt.Sprintf("%s/info?address=%s&skip_invalid=true&aux=status", c.baseURL, contract)
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
			if strings.EqualFold(ca.ContractAddress, contract) && ca.Platform.Coin.Name == chainName {
				c.cachedData.Set(c.getCacheKey(ca.Platform.Name, contract), v.ID, cache.DefaultExpiration)
				return v.ID, nil
			}
		}
	}
	return -1, fmt.Errorf("failed to get cmc id for contract: %s", contract)
}

func (c *CMCService) getCacheKey(chain, contract string) string {
	return fmt.Sprintf("%s_%s", chain, contract)
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
