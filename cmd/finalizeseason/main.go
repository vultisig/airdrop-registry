package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/config"
	"github.com/vultisig/airdrop-registry/internal/services"
)

// At the end of the season, we need to commit the points of all vaults and reset the current season points
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load config")
	}
	storage, err := services.NewStorage(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize storage")
	}
	defer func() {
		if err := storage.Close(); err != nil {
			logrus.WithError(err).Fatal("Failed to close storage")
		}
	}()
	var currentVaultId uint
	for {
		vaults, err := storage.GetVaultsWithPage(currentVaultId, 100)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to get vaults")
		}
		if len(vaults) == 0 {
			break
		}
		for _, vault := range vaults {
			if err := storage.CommitSeasonPoints(vault.ID); err != nil {
				logrus.WithError(err).Fatalf("Failed to commit season points for vault %d", vault.ID)
			}
		}
		currentVaultId = vaults[len(vaults)-1].ID
	}
}
