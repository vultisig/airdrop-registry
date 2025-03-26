package tokens

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

const (
	tokenProgramID = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	jsonRPCVersion = "2.0"
	rpcMethod      = "getTokenAccountsByOwner"
)

type (
	jsonRpcResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  result `json:"result"`
	}
	result struct {
		Value []accountEntry `json:"value"`
	}
	accountEntry struct {
		Account account `json:"account"`
		Pubkey  string  `json:"pubkey"`
	}
	account struct {
		Data struct {
			Parsed struct {
				Info struct {
					IsNative    bool        `json:"isNative"`
					Mint        string      `json:"mint"`
					Owner       string      `json:"owner"`
					State       string      `json:"state"`
					TokenAmount tokenAmount `json:"tokenAmount"`
				} `json:"info"`
				Type string `json:"type"`
			} `json:"parsed"`
			Program string `json:"program"`
		} `json:"data"`
		Owner     string `json:"owner"`
		RentEpoch int64  `json:"rentEpoch"`
	}
	tokenAmount struct {
		Amount         string  `json:"amount"`
		Decimals       int     `json:"decimals"`
		UIAmount       float64 `json:"uiAmount"`
		UIAmountString string  `json:"uiAmountString"`
	}
)
type solDiscoveryService struct {
	logger         *logrus.Logger
	baseAddress    string
	cmcIDService   *CMCIDService
	oneinchService *oneinchService
}

func NewSolDiscoveryService(cmcIDService *CMCIDService) autoDiscoveryService {
	return &solDiscoveryService{
		logger:       logrus.WithField("module", "sol_discovery_service").Logger,
		baseAddress:  "https://api.mainnet-beta.solana.com",
		cmcIDService: cmcIDService,
	}
}

func (s *solDiscoveryService) discover(address string, chain common.Chain) ([]models.CoinBase, error) {
	if address == "" {
		return nil, fmt.Errorf("empty address provided")
	}

	response, err := s.fetchTokenAccounts(address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token accounts: %w", err)
	}

	return s.processTokenAccounts(response, address, chain)
}

func (s *solDiscoveryService) fetchTokenAccounts(address string) (*jsonRpcResponse, error) {

	requestBody := map[string]any{
		"jsonrpc": jsonRPCVersion,
		"id":      1,
		"method":  rpcMethod,
		"params": []any{
			address,
			map[string]string{
				"programId": tokenProgramID,
			},
			map[string]string{
				"encoding": "jsonParsed",
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(s.baseAddress, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"error": err,
			"url":   s.baseAddress,
			"chain": common.Solana,
		}).Error("Failed to fetch account balances")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response jsonRpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		s.logger.WithError(err).Error("Failed to decode response")
		return nil, err
	}

	return &response, nil
}

func (s *solDiscoveryService) processTokenAccounts(response *jsonRpcResponse, address string, chain common.Chain) ([]models.CoinBase, error) {
	coins := make([]models.CoinBase, 0)
	for _, entry := range response.Result.Value {
		info := entry.Account.Data.Parsed.Info
		if info.TokenAmount.Amount == "0" {
			continue
		}

		cmcid, err := s.cmcIDService.GetCMCIDByContract("Solana", info.Mint)
		if err != nil {
			s.logger.WithError(err).WithField("contract", info.Mint).
				Warn("failed to get CMCID for contract")
			continue
		}

		coinBase := models.CoinBase{
			Address:         address,
			Balance:         info.TokenAmount.Amount,
			Chain:           chain,
			ContractAddress: info.Mint,
			CMCId:           cmcid,
		}
		coins = append(coins, coinBase)
	}
	return coins, nil
}

func (s *solDiscoveryService) search(coin models.CoinBase) (models.CoinBase, error) {
	coins, err := s.discover(coin.Address, coin.Chain)
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
