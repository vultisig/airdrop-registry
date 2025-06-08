package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/utils"
)

func (b *BalanceResolver) fetchERC721TokenBalance(chain common.Chain, contractAddress, address string) (float64, error) {
	if contractAddress == "" {
		return 0, fmt.Errorf("contract address cannot be empty")
	}
	if address == "" {
		return 0, fmt.Errorf("address cannot be empty")
	}
	baseUrl, err := b.getRpcUrlForChain(chain)
	if err != nil {
		return 0, fmt.Errorf("error getting rpc url for chain %s: %w", chain, err)
	}
	// Function signature hash of `balanceOf(address)` is `0x70a08231`
	functionSignature := "0x70a08231"
	// The wallet address is stripped of '0x', left-padded with zeros to 64 characters
	strippedWalletAddress := strings.TrimPrefix(address, "0x")
	paddedWalletAddress := fmt.Sprintf("%064s", strippedWalletAddress)

	// Concatenate function signature and padded wallet address
	data := strings.ToLower(functionSignature + paddedWalletAddress)

	// Create parameters array
	params := []interface{}{
		RpcParams{
			To:   contractAddress,
			Data: data,
		},
		"latest",
	}

	// Create RPC request
	rpcRequest := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		Params:  params,
		Id:      1,
	}

	// Convert RPC request to JSON
	requestBody, err := json.Marshal(rpcRequest)
	if err != nil {
		return 0, fmt.Errorf("error marshalling RPC request: %w", err)
	}
	// Send HTTP POST request
	resp, err := http.Post(baseUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer b.closer(resp.Body)

	if resp.StatusCode == http.StatusTooManyRequests {
		// rate limited, need to backoff and then retry
		return 0, ErrRateLimited
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance of address %s on %s: %s", address, chain, resp.Status)
	}
	// Parse response
	var rpcResponse RpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return 0, fmt.Errorf("error decoding RPC response: %w", err)
	}

	return utils.HexToFloat64(rpcResponse.Result, 0)
}
