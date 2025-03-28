package tokens

import (
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

type autoDiscoveryService interface {
	discover(address string, chain common.Chain) ([]models.CoinBase, error)
	//search(coin models.CoinBase) (models.CoinBase, error)
}
