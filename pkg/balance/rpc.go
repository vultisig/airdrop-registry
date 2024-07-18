package balance

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
