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

func GetBaseTokenByChain(chain string) (string, error) {
	switch chain {
	case "ethereum":
		return "ETH", nil
	case "avalanche":
		return "AVAX", nil
	case "bsc":
		return "BNB", nil
	case "base":
		return "ETH", nil
	case "blast":
		return "BLAST", nil
	case "optimism":
		return "ETH", nil
	case "matic":
		return "MATIC", nil
	case "zksync":
		return "ETH", nil
	case "bitcoin":
		return "BTC", nil
	case "mayachain":
		return "MAYA", nil
	case "thorchain":
		return "RUNE", nil
	default:
		return "", fmt.Errorf("unsupported chain: %s", chain)
	}
}
