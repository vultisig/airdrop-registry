package tokens

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/mr-tron/base58"
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

type (
	trcAccountResponse struct {
		Data    []trcAccount `json:"data"`
		Success bool         `json:"success"`
	}

	trcAccount struct {
		Address string              `json:"address"`
		Balance int64               `json:"balance"`
		Trc20   []map[string]string `json:"trc20"`
	}
)
type trcDiscoveryService struct {
	logger       *logrus.Logger
	tronBaseURL  string
	cmcIDService *CMCIDService
}

func NewTrcDiscoveryService(chain common.Chain, cmcIDService *CMCIDService) autoDiscoveryService {
	return &trcDiscoveryService{
		logger:       logrus.WithField("module", "trc_service").Logger,
		tronBaseURL:  "https://api.trongrid.io",
		cmcIDService: cmcIDService,
	}
}

func (trc *trcDiscoveryService) discover(address string, chain common.Chain) ([]models.CoinBase, error) {
	// Validate inputs
	if address == "" {
		return nil, fmt.Errorf("empty address provided")
	}
	if chain != common.Tron {
		return nil, fmt.Errorf("chain does not support")
	}
	escaped := url.PathEscape(address)
	accountURL := fmt.Sprintf("%s/v1/accounts/%s", trc.tronBaseURL, escaped)
	resp, err := http.Get(accountURL)
	if err != nil {
		trc.logger.WithError(err).Errorf("failed to fetch account from %s", accountURL)
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		trc.logger.Errorf("unexpected status code: %d for %s", resp.StatusCode, accountURL)
		return nil, fmt.Errorf("unexpecsted status code: %d", resp.StatusCode)
	}

	var accountResponse trcAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountResponse); err != nil {
		trc.logger.WithError(err).Error("failed to decode response")
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !accountResponse.Success {
		trc.logger.Warn("unsuccessful response from TRC API")
		return nil, fmt.Errorf("unsuccessful response from TRC API")
	}

	coins, err := trc.processAccounts(address, accountResponse.Data)
	if err != nil {
		trc.logger.Warn("failed to proccess accounts")
		return nil, fmt.Errorf("failed to proccess accounts: %w", err)
	}
	return coins, nil
}

type symbolDecimalModel struct {
	Result struct {
		Result bool `json:"result"`
	} `json:"result"`
	EnergyUsed     int      `json:"energy_used"`
	ConstantResult []string `json:"constant_result"`
	EnergyPenalty  int      `json:"energy_penalty"`
	Transaction    struct {
		Ret     []struct{} `json:"ret"`
		Visible bool       `json:"visible"`
		TxID    string     `json:"txID"`
		RawData struct {
			Contract []struct {
				Parameter struct {
					Value struct {
						Data            string `json:"data"`
						OwnerAddress    string `json:"owner_address"`
						ContractAddress string `json:"contract_address"`
					} `json:"value"`
					TypeURL string `json:"type_url"`
				} `json:"parameter"`
				Type string `json:"type"`
			} `json:"contract"`
			RefBlockBytes string `json:"ref_block_bytes"`
			RefBlockHash  string `json:"ref_block_hash"`
			Expiration    int64  `json:"expiration"`
			Timestamp     int64  `json:"timestamp"`
		} `json:"raw_data"`
		RawDataHex string `json:"raw_data_hex"`
	} `json:"transaction"`
}

func (trc *trcDiscoveryService) processAccounts(address string, accounts []trcAccount) ([]models.CoinBase, error) {
	coins := make([]models.CoinBase, 0)
	for _, account := range accounts {
		coinBases, err := trc.processTRC20Tokens(address, account.Trc20)
		if err != nil {
			trc.logger.WithError(err).Warn("error processing TRC20 tokens")
			continue
		}
		coins = append(coins, coinBases...)
	}
	return coins, nil
}

func (trc *trcDiscoveryService) processTRC20Tokens(address string, tokens []map[string]string) ([]models.CoinBase, error) {
	coins := make([]models.CoinBase, 0)
	for _, tokenMap := range tokens {
		for contract, balanceStr := range tokenMap {
			coin, err := trc.processToken(address, contract, balanceStr)
			if err != nil {
				trc.logger.WithError(err).WithField("contract", contract).Warn("error processing token")
				continue
			}
			if coin != nil {
				coins = append(coins, *coin)
			}
		}
	}
	return coins, nil
}

func (trc *trcDiscoveryService) processToken(address, contract, balanceStr string) (*models.CoinBase, error) {
	balance := new(big.Int)
	if _, ok := balance.SetString(balanceStr, 10); !ok {
		return nil, fmt.Errorf("invalid balance: %s", balanceStr)
	}

	if balance.Cmp(big.NewInt(0)) <= 0 {
		return nil, nil
	}

	cmcid, err := trc.cmcIDService.GetCMCIDByContract("TRON", contract)
	if err != nil {
		return nil, fmt.Errorf("failed to get CMCID: %w", err)
	}

	symbol, err := trc.fetchTokenData(address, contract, "symbol()")
	if err != nil {
		return nil, fmt.Errorf("failed to get symbol: %w", err)
	}

	decimalsHex, err := trc.fetchTokenData(address, contract, "decimals()")
	if err != nil {
		return nil, fmt.Errorf("failed to get decimals: %w", err)
	}

	decimal, err := strconv.ParseInt(decimalsHex, 16, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decimals: %w", err)
	}

	return &models.CoinBase{
		Ticker:          symbol,
		Address:         address,
		Balance:         balanceStr,
		CMCId:           cmcid,
		Chain:           common.Tron,
		ContractAddress: contract,
		Decimals:        int(decimal),
	}, nil
}

func (trc *trcDiscoveryService) fetchTokenData(address, contract, selector string) (string, error) {
	hexContract, err := DecodeBase58ToHex(contract)
	if err != nil {
		return "", fmt.Errorf("failed to decode contract hex: %w", err)
	}

	hexAddress, err := DecodeBase58ToHex(address)
	if err != nil {
		return "", fmt.Errorf("failed to decode contract hex: %w", err)
	}
	url := fmt.Sprintf("%s/wallet/triggerconstantcontract", trc.tronBaseURL)
	payload := fmt.Sprintf(`{"contract_address": "%s","function_selector": "%s","owner_address": "%s"}`, hexContract, selector, hexAddress)

	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("failed to fetch data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response symbolDecimalModel
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.ConstantResult) == 0 {
		return "", fmt.Errorf("no data returned")
	}

	result := response.ConstantResult[0]
	if selector == "symbol()" {
		if len(result) < 128 {
			return "", fmt.Errorf("invalid symbol data length")
		}
		symbolBytes, err := hexToBytes(result[128:])
		if err != nil {
			return "", fmt.Errorf("failed to decode symbol hex: %w", err)
		}
		return string(bytes.TrimRight(symbolBytes, "\x00")), nil
	}
	return result, nil
}

func (trc *trcDiscoveryService) search(coin models.CoinBase) (models.CoinBase, error) {
	coins, err := trc.discover(coin.Address, coin.Chain)
	if err != nil {
		return models.CoinBase{}, fmt.Errorf("failed to discover tokens: %w", err)
	}

	for _, c := range coins {
		if c.ContractAddress == coin.ContractAddress {
			return c, nil
		}
	}
	return models.CoinBase{}, fmt.Errorf("token not found")
}

func DecodeBase58ToHex(base58Address string) (string, error) {
	rawBytes, err := base58.Decode(base58Address)
	if err != nil {
		return "", err
	}

	payload := rawBytes[:len(rawBytes)-4]
	hexAddress := fmt.Sprintf("%X", payload)
	return hexAddress, nil
}

func HexToBase58(hexAddress string) (string, error) {
	rawBytes, err := hex.DecodeString(hexAddress)
	if err != nil {
		return "", err
	}

	hash1 := sha256.Sum256(rawBytes)
	hash2 := sha256.Sum256(hash1[:])
	checksum := hash2[:4]

	addressBytes := append(rawBytes, checksum...)
	base58Address := base58.Encode(addressBytes)
	return base58Address, nil
}

func hexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	return hex.DecodeString(hexStr)
}
