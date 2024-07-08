package balance

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vultisig/airdrop-registry/pkg/utils"
)

func FetchEvmBalanceOfAddress(chain, address string) (float64, error) {
	rpcUrl := getRpcUrlForChain(chain)
	if rpcUrl == "" {
		return 0, fmt.Errorf("unsupported EVM chain: %s", chain)
	}

	payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBalance","params":["%s", "latest"],"id":1}`, address)
	response, err := http.Post(rpcUrl, "application/json", strings.NewReader(payload))
	if err != nil {
		return 0, fmt.Errorf("error fetching balance of address %s on %s: %v", address, chain, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response body: %v", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	balanceHex := data["result"].(string)
	balance, err := utils.HexToFloat64(balanceHex)
	if err != nil {
		return 0, fmt.Errorf("error converting balance to float: %v", err)
	}

	return balance, nil
}

func getRpcUrlForChain(chain string) string {
	switch chain {
	case "ethereum":
		return "https://ethereum-rpc.publicnode.com"
	case "avalanche":
		return "https://api.avax.network/ext/bc/C/rpc"
	case "bsc":
		return "https://bsc-dataseed.binance.org/"
	case "eth_base":
		return "https://mainnet.base.org"
	case "eth_blast":
		return "https://rpc.ankr.com/eth"
	case "optimism":
		return "https://mainnet.optimism.io"
	case "matic":
		return "https://polygon-rpc.com"
	case "zksync":
		return "https://mainnet.era.zksync.io"
	default:
		return ""
	}
}
