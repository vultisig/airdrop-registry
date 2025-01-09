package services

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/liquidity"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func TestPointService(t *testing.T) {
	priceResolver, err := NewPriceResolver(&config.Config{})
	if err != nil {
		t.Errorf("Failed to create price resolver: %v", err)
		t.FailNow()
	}
	pointService := PointWorker{
		logger:          logrus.WithField("module", "point_service").Logger,
		storage:         nil,
		priceResolver:   priceResolver,
		balanceResolver: nil,
		lpResolver:      liquidity.NewLiquidtyPositionResolver(),
		saverResolver:   liquidity.NewSaverPositionResolver(),
	}
	vaultAddress := models.NewVaultAddress(1064)
	vaultAddress.SetAddress(common.Ethereum, "0x562f334890C717f31bAB4c1197C67619FbD0eAFc")
	vaultAddress.SetAddress(common.Bitcoin, "bc1quqf7l4kvjtwswu29n6nd4szr54khpzmh4yp3cn")
	vaultAddress.SetAddress(common.THORChain, "thor1rj7eqmqmeyrvkmlmrded42609gedtqnuwsmhww")
	vaultAddress.SetAddress(common.Solana, "CbSjseduYqKiavFxvdeRVH6DBv9Fz4rd59BLAFJz8J9Q")
	vaultAddress.SetAddress(common.BscChain, "0x562f334890C717f31bAB4c1197C67619FbD0eAFc")
	for i := 0; i < 1; i++ {
		newLPValue, err := pointService.fetchPosition(vaultAddress)
		if err != nil {
			t.Error(err)
		}
		t.Logf("New LP Value: %v", newLPValue)
	}
}
