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
	var rpcResp RpcSolanaResp
	if err := json.NewDecoder(response.Body).Decode(&rpcResp); err != nil {
		return 0, fmt.Errorf("error decoding response: %v", err)
	}
	return rpcResp.Result.Value / 1000000000, nil
}
