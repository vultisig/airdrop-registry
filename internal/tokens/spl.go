package tokens

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/internal/utils"
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

func (s *splDiscoveryService) Search(coin models.CoinBase) (models.CoinBase, error) {
	cmcId, err := s.cmcService.GetCMCIDByContract(coin.Chain.String(), coin.ContractAddress)
	if err != nil {
		s.logger.WithError(err).WithField("contract", coin.ContractAddress).
			Warn("failed to get CMCID for contract")
		return models.CoinBase{}, err
	}
	decimal, err := s.getCoinDecimal(coin.ContractAddress)
	if err != nil {
		s.logger.WithError(err).WithField("contract", coin.ContractAddress).
			Warn("failed to get decimal for contract")
		return models.CoinBase{}, err
	}
	coin.CMCId = cmcId
	coin.Decimals = decimal
	return coin, nil
}

func (s *splDiscoveryService) getCoinDecimal(address string) (int, error) {
	parmas := []interface{}{
		address,
		map[string]string{
			"encoding": "jsonParsed",
		},
	}
	rpcRequest := utils.NewJsonRPCRequest("getAccountInfo", parmas, 1)
	buf, err := json.Marshal(rpcRequest)
	if err != nil {
		return 0, fmt.Errorf("error marshalling RPC request: %w", err)
	}
	resp, err := http.Post(s.baseAddress, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s: %w", address, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		// rate limited, need to backoff and then retry
		return 0, fmt.Errorf("rate limited while fetching balance of address %s", address)
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance of address %s: %s", address, resp.Status)
	}
	var result SPLTokenInfoResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}
	return result.Result.Value.Data.Parsed.Info.Decimals, nil
}

type SPLTokenInfoResp struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Context struct {
			APIVersion string `json:"apiVersion"`
			Slot       int    `json:"slot"`
		} `json:"context"`
		Value struct {
			Data struct {
				Parsed struct {
					Info struct {
						Decimals        int    `json:"decimals"`
						FreezeAuthority string `json:"freezeAuthority"`
						IsInitialized   bool   `json:"isInitialized"`
						MintAuthority   string `json:"mintAuthority"`
						Supply          string `json:"supply"`
					} `json:"info"`
					Type string `json:"type"`
				} `json:"parsed"`
				Program string `json:"program"`
				Space   int    `json:"space"`
			} `json:"data"`
			Executable bool   `json:"executable"`
			Lamports   int64  `json:"lamports"`
			Owner      string `json:"owner"`
			RentEpoch  int64  `json:"rentEpoch"`
			Space      int    `json:"space"`
		} `json:"value"`
	} `json:"result"`
	ID int `json:"id"`
}
