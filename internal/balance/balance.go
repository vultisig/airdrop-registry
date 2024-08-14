package balance

import (
	"github.com/sirupsen/logrus"
)

const vultisigApiProxy = "https://api.vultisig.com"

// BalanceResolver is to fetch address balances
type BalanceResolver struct {
	logger *logrus.Logger
}

func NewBalanceResolver() (*BalanceResolver, error) {
	return &BalanceResolver{
		logger: logrus.WithField("module", "balance_resolver").Logger,
	}, nil
}

//
//func FetchBalanceOfAddress(chain, address string) (float64, error) {
//	switch chain {
//	case "bitcoin":
//		return FetchBitcoinBalanceOfAddress(address)
//	case "ethereum", "avalanche", "bsc", "eth_base", "eth_blast", "optimism", "matic", "zksync":
//		// @TEST
//		if chain == "ethereum" {
//			address = "0xaA11EA95475341c4dDb83aF141B01e52500c23d6"
//		}
//		return FetchEvmBalanceOfAddress(chain, address)
//	case "mayachain":
//		return FetchMayachainBalanceOfAddress(address)
//	case "thorchain":
//		return FetchThorchainBalanceOfAddress(address)
//	case "polkadot":
//		return FetchPolkadotBalanceOfAddress(address)
//	case "sui":
//		return FetchSuiBalanceOfAddress(address)
//	case "solana":
//		return FetchSolanaBalanceOfAddress(address)
//	default:
//		return 0, fmt.Errorf("unsupported chain: %s", chain)
//	}
//}
//
//func GetBaseTokenByChain(chain string) (string, error) {
//	switch chain {
//	case "ethereum":
//		return "ETH", nil
//	case "avalanche":
//		return "AVAX", nil
//	case "bsc":
//		return "BNB", nil
//	case "base":
//		return "ETH", nil
//	case "blast":
//		return "BLAST", nil
//	case "optimism":
//		return "ETH", nil
//	case "matic":
//		return "MATIC", nil
//	case "zksync":
//		return "ETH", nil
//	case "bitcoin":
//		return "BTC", nil
//	case "litecoin":
//		return "LTC", nil
//	case "dash":
//		return "DASH", nil
//	case "bitcoincash":
//		return "BCH", nil
//	case "dogecoin":
//		return "DOGE", nil
//	case "mayachain":
//		return "MAYA", nil
//	case "thorchain":
//		return "RUNE", nil
//	case "polkadot":
//		return "DOT", nil
//	case "sui":
//		return "SUI", nil
//	case "solana":
//		return "SOL", nil
//	default:
//		return "", fmt.Errorf("unsupported chain: %s", chain)
//	}
//}
