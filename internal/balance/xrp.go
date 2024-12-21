package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (b *BalanceResolver) FetchXRPBalanceOfAddress(address string) (float64, error) {
	rpcUrl := b.xrpBalanceBaseAddress
	// Create parameters array
	params := []interface{}{
		map[string]interface{}{
			"account":      address,
			"ledger_index": "current",
			"queue":        true,
		},
	}

	// Create RPC request
	rpcReq := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "account_info",
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
	type RpcXRPResp struct {
		Result struct {
			AccountData struct {
				Balance int64 `json:"Balance,string"`
			} `json:"account_data"`
		} `json:"result"`
	}
	var rpcResp RpcXRPResp
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	return float64(rpcResp.Result.AccountData.Balance) / 1e6, nil
}
