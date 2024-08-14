package balance

func GetChainIDByChain(chain string) int {
	switch chain {
	case "ethereum":
		return 1 // Ethereum
	case "avalanche":
		return 43114 // Avalanche
	case "bsc":
		return 56 // BSC
	case "eth_base":
		return 8453 // ETH Base
	case "eth_blast":
		return 100 // ETH Blast
	case "optimism":
		return 10 // Optimism
	case "matic":
		return 137 // Matic
	case "zksync":
		return 324 // zkSync
	case "sui":
		return 1000 // Sui
	default:
		return 0 // Unknown chain
	}
}
