package balance

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

const (
	vultisigApiProxy = "https://api.vultisig.com"
	maxRetries       = 10
	initialBackoff   = time.Second
)

// BalanceResolver is to fetch address balances
type BalanceResolver struct {
	logger                 *logrus.Logger
	thorchainBondProviders *sync.Map
	thorchainRuneProviders *sync.Map
	thornodeBaseAddress    string
	tonBalanceBaseAddress  string
	tronBalanceBaseAddress string
	xrpBalanceBaseAddress  string
	whitelistNFTCollection []models.NFTCollection
	whiteListSPLToken      map[string]string
	whiteListTRC20Token    map[string]int
}

func NewBalanceResolver() (*BalanceResolver, error) {
	return &BalanceResolver{
		logger:                 logrus.WithField("module", "balance_resolver").Logger,
		thorchainBondProviders: &sync.Map{},
		thorchainRuneProviders: &sync.Map{},
		thornodeBaseAddress:    "https://thornode.ninerealms.com",
		tonBalanceBaseAddress:  "https://api.vultisig.com/ton/v3/addressInformation",
		tronBalanceBaseAddress: "https://api.trongrid.io",
		xrpBalanceBaseAddress:  "https://xrplcluster.com",
		whitelistNFTCollection: []models.NFTCollection{
			{
				Chain:             common.Ethereum,
				CollectionAddress: "0xa98b29a8f5a247802149c268ecf860b8308b7291",
				CollectionSlug:    "thorguards",
			},
		},
		whiteListSPLToken: map[string]string{
			"DEf93bSt8dx58gDFCcz4CwbjYZzjwaRBYAciJYLfdCA9": "KWEEN",
			"rndrizKT3MK1iimdxRdWabcF7Zg7AR5T4nud4EkHBof":  "RENDER",
			"EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v": "USDC",
			"Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB": "USDT",
			"JUPyiwrYJFskUPiHa7hkeR8VUtAeFoSYbKedZNsDvCN":  "JUP",
			"FgWto1nfArQTpg3o74sYkti753caPfHNXHG8CkedDpMg": "DORITO",
		},
		whiteListTRC20Token: map[string]int{
			"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t": 6, // USDT
		},
	}, nil
}

func (b *BalanceResolver) GetBalanceWithRetry(coin models.CoinDBModel) (float64, error) {
	var balance float64
	var err error

	for i := 0; i < maxRetries; i++ {
		balance, err = b.GetBalance(coin)
		if err == nil {
			return balance, nil
		}

		if !errors.Is(err, ErrRateLimited) {
			return 0, err
		}

		backoffDuration := initialBackoff * time.Duration(i)
		b.logger.Warnf("Rate limited. Retrying in %s...", backoffDuration)
		time.Sleep(backoffDuration)
	}

	return 0, fmt.Errorf("failed to get balance after %d retries: %w", maxRetries, err)
}

func (b *BalanceResolver) GetBalance(coin models.CoinDBModel) (float64, error) {
	switch coin.Chain {
	case common.Bitcoin, common.BitcoinCash, common.Litecoin, common.Dogecoin, common.Dash:
		balance, _, err := b.FetchUtxoBalanceOfAddress(coin.Address, coin.Chain)
		return balance, err
	case common.Arbitrum, common.Ethereum, common.Zksync, common.Optimism, common.Polygon, common.BscChain, common.Avalanche, common.Base, common.Blast, common.CronosChain:
		if coin.ContractAddress != "" {
			for _, nft := range b.whitelistNFTCollection {
				if strings.EqualFold(coin.ContractAddress, nft.CollectionAddress) {
					return b.fetchERC721TokenBalance(coin.Chain, coin.ContractAddress, coin.Address)
				}
			}
			return b.fetchERC20TokenBalance(coin.Chain, coin.ContractAddress, coin.Address, int64(coin.Decimals))
		} else {
			return b.FetchEvmBalanceOfAddress(coin.Chain, coin.Address)
		}
	case common.THORChain:
		return b.FetchThorchainBalanceOfAddress(coin.Address)
	case common.MayaChain:
		if strings.EqualFold(coin.Ticker, "maya") {
			return b.FetchMayachainMayaBalanceOfAddress(coin.Address)
		} else if strings.EqualFold(coin.Ticker, "cacao") {
			return b.FetchMayachainCacoBalanceOfAddress(coin.Address)
		}
	case common.GaiaChain:
		return b.FetchCosmosBalanceOfAddress(coin.Address)
	case common.Dydx:
		return b.FetchDydxBalanceOfAddress(coin.Address)
	case common.Terra:
		return b.FetchTerraBalanceOfAddress(coin.Address)
	case common.TerraClassic:
		return b.FetchTerraClassicBalanceOfAddress(coin.Address)
	case common.Noble:
		if strings.EqualFold(coin.Ticker, "USDC") { //  We only support USDC on Noble for now
			return b.FetchNobleBalanceOfAddress(coin.Address)
		}
	case common.Kujira:
		var totalBalance float64
		balanceKujira, errK := b.FetchKujiraBalanceOfAddress(coin.Address)
		if errK == nil {
			totalBalance += balanceKujira
		}
		balanceRkujira, errR := b.FetchRkujiraBalanceOfAddress(coin.Address)
		if errR == nil {
			totalBalance += balanceRkujira
		}
		return totalBalance, nil
	case common.Osmosis:
		return b.FetchOsmosisBalanceOfAddress(coin.Address)
	case common.Akash:
		return b.FetchAkashBalanceOfAddress(coin.Address)
	case common.Solana:
		//ignore none native coins (spl tokens)
		if coin.ContractAddress == "" {
			return b.FetchSolanaBalanceOfAddress(coin.Address)
		} else {
			for addr, _ := range b.whiteListSPLToken {
				if strings.EqualFold(coin.ContractAddress, addr) {
					return b.FetchSPLBalanceOfAddress(coin.Address, coin.ContractAddress)
				}
			}
			return 0, nil
		}
	case common.Polkadot:
		return b.FetchPolkadotBalanceOfAddress(coin.Address)
	case common.Sui:
		return b.FetchSuiBalanceOfAddress(coin.Address)
	case common.Ton:
		return b.FetchTonBalanceOfAddress(coin.Address)
	case common.XRP:
		return b.FetchXRPBalanceOfAddress(coin.Address)
	case common.Tron:
		if coin.ContractAddress == "" { // TRX token
			return b.FetchTronBalanceOfAddress(coin.Address, "", 6)
		} else {
			for addr, decimal := range b.whiteListTRC20Token {
				if coin.ContractAddress == addr {
					return b.FetchTronBalanceOfAddress(coin.Address, coin.ContractAddress, decimal)
				}
			}
			return 0, nil
		}
	default:
		return 0, fmt.Errorf("chain: %s doesn't support", coin.Chain)
	}
	return 0, nil
}
