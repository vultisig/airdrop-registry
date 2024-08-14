package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	//case common.Sui:
	//	return "https://sui-rpc.publicnode.com"
	default:
		return "", fmt.Errorf("chain: %s doesn't support", chain)
	}
}

const EVM_ETH_BALANCE_TEMPLATE = `
{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": [
        "%s",
        "latest"
    ],
    "id": 1
}
`

func (b *BalanceResolver) FetchEvmBalanceOfAddress(chain common.Chain, address string) (float64, error) {
	rpcUrl, err := b.getRpcUrlForChain(chain)
	if err != nil {
		return 0, fmt.Errorf("error getting rpc url for chain %s: %w", chain, err)
	}

	payload := fmt.Sprintf(EVM_ETH_BALANCE_TEMPLATE, address)
	resp, err := http.Post(rpcUrl, "application/json", strings.NewReader(payload))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on %s: %w", address, chain, err)
	}
	defer b.closer(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	balanceHex := data["result"].(string)
	balance, err := utils.HexToFloat64(balanceHex, 18)
	if err != nil {
		return 0, fmt.Errorf("error converting balance to float: %v", err)
	}

	return balance, nil
}
