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
	tokenProgramID     = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	jsonRPCVersion     = "2.0"
	discoveryRpcMethod = "getTokenAccountsByOwner"
)

type splDiscoveryService struct {
	logger      *logrus.Logger
	baseAddress string
	cmcService  *CMCService
}

func NewSPLDiscoveryService(cmcService *CMCService) AutoDiscoveryService {
	return &splDiscoveryService{
		logger:      logrus.WithField("module", "sol_discovery_service").Logger,
		baseAddress: "https://api.vultisig.com/solana/",
		cmcService:  cmcService,
	}
}

func (s *splDiscoveryService) Discover(address string, chain common.Chain) ([]models.CoinBase, error) {
	if address == "" {
		return nil, fmt.Errorf("empty address provided")
	}

	tokens, err := s.fetchTokenAccounts(address)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch token accounts: %w", err)
	}

	return tokens, nil
}

func (s *splDiscoveryService) fetchTokenAccounts(address string) ([]models.CoinBase, error) {

	requestBody := map[string]any{
		"jsonrpc": jsonRPCVersion,
		"id":      1,
		"method":  discoveryRpcMethod,
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

	coins := make([]models.CoinBase, 0)
	for _, entry := range response.Result.Value {
		info := entry.Account.Data.Parsed.Info
		if info.TokenAmount.Amount == "0" {
			continue
		}

		cmcid, err := s.cmcService.GetCMCIDByContract("Solana", info.Mint)
		if err != nil {
			s.logger.WithError(err).WithField("contract", info.Mint).
				Warn("failed to get CMCID for contract")
			continue
		}

		coinBase := models.CoinBase{
			Address:         address,
			Balance:         info.TokenAmount.Amount,
			Chain:           common.Solana,
			ContractAddress: info.Mint,
			CMCId:           cmcid,
		}
		coins = append(coins, coinBase)
	}
	return coins, nil
}

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
