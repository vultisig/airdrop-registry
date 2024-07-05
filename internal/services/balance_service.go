package services

import (
	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func SaveBalance(balance *models.Balance) error {
	return db.DB.Create(balance).Error
}

func GetBalancesByVaultID(vaultID uint) ([]models.Balance, error) {
	var balances []models.Balance
	err := db.DB.Where("vault_id = ?", vaultID).Find(&balances).Error
	return balances, err
}
