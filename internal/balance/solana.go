package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RpcSolanaResp struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  struct {
		Value float64 `json:"value"`
	} `json:"result"`
}

type RpcSplResp struct {
	Result struct {
		Value []struct {
			Account struct {
				Data struct {
					Parsed struct {
						Info struct {
							IsNative    bool   `json:"isNative"`
							Mint        string `json:"mint"`
							Owner       string `json:"owner"`
							State       string `json:"state"`
							TokenAmount struct {
								Amount         string  `json:"amount"`
								Decimals       int     `json:"decimals"`
								UIAmount       float64 `json:"uiAmount"`
								UIAmountString string  `json:"uiAmountString"`
							} `json:"tokenAmount"`
						} `json:"info"`
						Type string `json:"type"`
					} `json:"parsed"`
				} `json:"data"`
			} `json:"account"`
		} `json:"value"`
	} `json:"result"`
}

func (b *BalanceResolver) FetchSolanaBalanceOfAddress(address string) (float64, error) {
	// Create parameters array
	params := []interface{}{
		address,
	}

	// Create RPC request
	rpcReq := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "getBalance",
		Params:  params,
		Id:      1,
	}
	// Convert RPC request to JSON
	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return 0, fmt.Errorf("error marshalling RPC request: %w", err)
	}
	response, err := http.Post("https://api.mainnet-beta.solana.com", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on Solana: %w", address, err)
	}
	defer b.closer(response.Body)
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance of address %s on Solana: %s", address, response.Status)
	}
	if response.StatusCode == http.StatusTooManyRequests {
		return 0, ErrRateLimited
	}
	var rpcResp RpcSolanaResp
	if err := json.NewDecoder(response.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}
	return rpcResp.Result.Value / 1000000000, nil
}

func (b *BalanceResolver) FetchSPLBalanceOfAddress(vaultAddress, contractAdderss string) (float64, error) {
	// Create parameters array
	params := []interface{}{
		vaultAddress,
		map[string]interface{}{
			"mint": contractAdderss,
		},
		map[string]interface{}{
			"encoding": "jsonParsed",
		},
	}

	// Create RPC request
	rpcReq := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "getTokenAccountsByOwner",
		Params:  params,
		Id:      1,
	}
	// Convert RPC request to JSON
	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return 0, fmt.Errorf("error marshalling RPC request: %w", err)
	}
	response, err := http.Post("https://api.mainnet-beta.solana.com", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("error fetching spl balance  %s of address %s on Solana: %w", contractAdderss, vaultAddress, err)
	}
	defer b.closer(response.Body)
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching spl balance %s of address %s on Solana: %s", contractAdderss, vaultAddress, response.Status)
	}
	if response.StatusCode == http.StatusTooManyRequests {
		return 0, ErrRateLimited
	}
	var rpcResp RpcSplResp
	if err := json.NewDecoder(response.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}
	for _, v := range rpcResp.Result.Value {
		if v.Account.Data.Parsed.Info.Mint == contractAdderss {
			return v.Account.Data.Parsed.Info.TokenAmount.UIAmount, nil
		}
	}
	return 0, nil
}
