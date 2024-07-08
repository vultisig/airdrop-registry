package balance

import "fmt"

func FetchBalanceOfAddress(chain, address string) (float64, error) {
	switch chain {
	case "bitcoin":
		return FetchBitcoinBalanceOfAddress(address)
	case "ethereum", "avalanche", "bsc", "eth_base", "eth_blast", "optimism", "matic", "zksync":
		return FetchEvmBalanceOfAddress(chain, address)
	case "mayachain":
		return FetchMayachainBalanceOfAddress(address)
	case "thorchain":
		return FetchThorchainBalanceOfAddress(address)
	default:
		return 0, fmt.Errorf("unsupported chain: %s", chain)
	}
}
