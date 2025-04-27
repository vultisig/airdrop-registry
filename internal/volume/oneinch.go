package volume

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/utils"
)

type oneInchVolumeTracker struct {
	etherscanbaseUrl string
	etherscanApiKey  string
	ethplorerBaseUrl string
	ethplorerApiKey  string
	logger           *logrus.Logger
}

func NewOneInchVolumeTracker(etherscanApiKey, ethplorerApiKey string) IVolumeTracker {
	return &oneInchVolumeTracker{
		etherscanbaseUrl: "https://api.etherscan.io",
		ethplorerBaseUrl: "https://api.ethplorer.io",
		etherscanApiKey:  etherscanApiKey,
		ethplorerApiKey:  ethplorerApiKey,
		logger:           logrus.WithField("module", "oneInch_volume_tracker").Logger,
	}
}
func (o *oneInchVolumeTracker) SafeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		o.logger.Error(err)
	}
}

// TODO: add from/to filter to etherscan request
func (o *oneInchVolumeTracker) FetchVolume(from, to int64, affiliate string) (map[string]float64, error) {
	res := make(map[string]float64)
	//ignore invalid affiliate
	if !o.isValidAffiliate(affiliate) {
		return res, nil
	}
	// #TODO check api for from & to parameters
	url := fmt.Sprintf("%s/v2/api?chainid=1&module=account&action=txlistinternal&address=%s&apikey=%s", o.etherscanbaseUrl, affiliate, o.etherscanApiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer o.SafeClose(resp.Body)
	var etherScanResponse, etherScanFilteredResponse etherScanResponse
	if err := json.NewDecoder(resp.Body).Decode(&etherScanResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	if etherScanResponse.Status != "1" {
		return nil, fmt.Errorf("error in etherscan response: %s", etherScanResponse.Message)
	}
	for _, item := range etherScanResponse.Result {
		if from <= item.TimeStamp && item.TimeStamp < to {
			etherScanFilteredResponse.Result = append(etherScanFilteredResponse.Result, item)
		}
	}
	txHashes := make([]string, 0)
	for _, item := range etherScanFilteredResponse.Result {
		if item.IsError != 0 {
			continue
		}
		txHashes = append(txHashes, item.Hash)
	}

	for _, tx := range txHashes {
		url = fmt.Sprintf("%s/getTxInfo/%s?apiKey=%s", o.ethplorerBaseUrl, tx, o.ethplorerApiKey)
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("error making GET request: %w", err)
		}
		defer o.SafeClose(resp.Body)
		var ethplorerTx ethplorerVolumeModel
		if err := json.NewDecoder(resp.Body).Decode(&ethplorerTx); err != nil {
			return nil, fmt.Errorf("error decoding response: %w", err)
		}
		if len(ethplorerTx.Operations) == 0 {
			continue
		}
		lastOperation := ethplorerTx.Operations[len(ethplorerTx.Operations)-1]
		price, ok := lastOperation.TokenInfo.Price.(map[string]any)
		if ok {
			var priceData ethplorerPrice
			err = mapToStruct(price, &priceData)
			if err != nil {
				return nil, fmt.Errorf("error converting map to struct: %w", err)
			}
			bigIntValue, ok := new(big.Int).SetString(lastOperation.Value, 10)
			if !ok {
				return nil, fmt.Errorf("error converting string to big int: %s", lastOperation.Value)
			}
			floatValue := new(big.Float).SetInt(bigIntValue)
			scale := new(big.Float).SetFloat64(math.Pow10(lastOperation.TokenInfo.Decimals))
			amount, _ := new(big.Float).Quo(floatValue, scale).Float64()
			res[lastOperation.To] += amount * priceData.Rate
		} else {
			res[lastOperation.To] = 0
		}
	}
	return res, nil
}

// Affilaite should be eth address
func (o *oneInchVolumeTracker) isValidAffiliate(affiliate string) bool {
	return utils.IsETHAddress(affiliate)
}

func mapToStruct(m map[string]any, result any) error {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonStr, result)
}

type etherScanResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Result  []etherscanResult `json:"result"`
}
type etherscanResult struct {
	Hash      string `json:"hash"`
	TimeStamp int64  `json:"timeStamp,string"`
	IsError   int    `json:"isError,string"`
}
type ethplorerPrice struct {
	Rate float64 `json:"rate"`
}
type ethplorerTokenInfo struct {
	Decimals int `json:"decimals,string"`
	Price    any `json:"price"`
}
type ethplorerOperation struct {
	Value     string             `json:"value"`
	To        string             `json:"to"`
	TokenInfo ethplorerTokenInfo `json:"tokenInfo"`
}
type ethplorerVolumeModel struct {
	Operations []ethplorerOperation `json:"operations"`
}
