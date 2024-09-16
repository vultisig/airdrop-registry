package balance

import (
	"errors"
	"fmt"
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
}

func NewBalanceResolver() (*BalanceResolver, error) {
	return &BalanceResolver{
		logger:                 logrus.WithField("module", "balance_resolver").Logger,
		thorchainBondProviders: &sync.Map{},
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
			return b.fetchERC20TokenBalance(coin.Chain, coin.ContractAddress, coin.Address, int64(coin.Decimals))
		} else {
			return b.FetchEvmBalanceOfAddress(coin.Chain, coin.Address)
		}
	case common.THORChain:
		return b.FetchThorchainBalanceOfAddress(coin.Address)
	case common.MayaChain:
		return b.FetchMayachainBalanceOfAddress(coin.Address)
	case common.GaiaChain:
		return b.FetchCosmosBalanceOfAddress(coin.Address)
	case common.Dydx:
		return b.FetchDydxBalanceOfAddress(coin.Address)
	case common.Kujira:
		return b.FetchKujiraBalanceOfAddress(coin.Address)
	case common.Solana:
		return b.FetchSolanaBalanceOfAddress(coin.Address)
	case common.Polkadot:
		return b.FetchPolkadotBalanceOfAddress(coin.Address)
	case common.Sui:
		return b.FetchSuiBalanceOfAddress(coin.Address)
	default:
		return 0, fmt.Errorf("chain: %s doesn't support", coin.Chain)
	}
}
