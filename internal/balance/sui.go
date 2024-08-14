package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (b *BalanceResolver) FetchSuiBalanceOfAddress(address string) (float64, error) {
	rpcUrl := "https://sui-rpc.publicnode.com"
	// Create parameters array
	params := []interface{}{
		address,
	}

	// Create RPC request
	rpcReq := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "suix_getBalance",
		Params:  params,
		Id:      1,
	}

	// Convert RPC request to JSON
	reqBody, err := json.Marshal(rpcReq)
	if err != nil {
		return 0, fmt.Errorf("error marshalling RPC request: %w", err)
	}

	resp, err := http.Post(rpcUrl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on SUI: %w", address, err)
	}
	defer b.closer(resp.Body)
	if resp.StatusCode == http.StatusTooManyRequests {
		// rate limited, need to backoff and then retry
		return 0, ErrRateLimited
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance of address %s on SUI: %s", address, resp.Status)
	}
	type RpcSuiResp struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  struct {
			TotalBalance string `json:"totalBalance"`
		} `json:"result"`
	}
	var rpcResp RpcSuiResp
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	balance, err := strconv.ParseFloat(rpcResp.Result.TotalBalance, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting balance to float: %w", err)
	}

	balance = balance / 1e9

	return balance, nil
}
