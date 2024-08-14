package balance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/utils"
)

func (b *BalanceResolver) getRpcUrlForChain(chain common.Chain) (string, error) {
	switch chain {
	case common.Ethereum:
		return "https://ethereum-rpc.publicnode.com", nil
	case common.Avalanche:
		return "https://avalanche-c-chain-rpc.publicnode.com", nil
	case common.BscChain:
		return "https://bsc-rpc.publicnode.com", nil
	case common.Base:
		return "https://base-rpc.publicnode.com", nil
	case common.Blast:
		return "https://rpc.ankr.com/blast", nil
	case common.Optimism:
		return "https://optimism-rpc.publicnode.com", nil
	case common.Polygon:
		return "https://polygon-bor-rpc.publicnode.com", nil
	case common.Zksync:
		return "https://mainnet.era.zksync.io", nil
	default:
		return "", fmt.Errorf("chain: %s doesn't support", chain)
	}
}

func (b *BalanceResolver) FetchEvmBalanceOfAddress(chain common.Chain, address string) (float64, error) {
	rpcUrl, err := b.getRpcUrlForChain(chain)
	if err != nil {
		return 0, fmt.Errorf("error getting rpc url for chain %s: %w", chain, err)
	}
	// Create parameters array
	params := []interface{}{
		address,
		"latest",
	}

	// Create RPC request
	rpcRequest := RpcRequest{
		Jsonrpc: "2.0",
		Method:  "eth_getBalance",
		Params:  params,
		Id:      1,
	}
	buf, err := json.Marshal(rpcRequest)
	if err != nil {
		return 0, fmt.Errorf("error marshalling RPC request: %w", err)
	}
	resp, err := http.Post(rpcUrl, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on %s: %w", address, chain, err)
	}
	defer b.closer(resp.Body)
	if resp.StatusCode == http.StatusTooManyRequests {
		// rate limited, need to backoff and then retry
		return 0, ErrRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("error fetching balance of address %s on %s: %s", address, chain, resp.Status)
	}
	type EthBalanceResult struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  string `json:"result"`
	}
	var result EthBalanceResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("error decoding response: %w", err)
	}

	balance, err := utils.HexToFloat64(result.Result, 18)
	if err != nil {
		return 0, fmt.Errorf("error converting balance to float: %w", err)
	}

	return balance, nil
}
