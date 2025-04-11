package tokens

import (
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

type AutoDiscoveryService interface {
	Discover(address string, chain common.Chain) ([]models.CoinBase, error)
	Search(coin models.CoinBase) (models.CoinBase, error)
}
