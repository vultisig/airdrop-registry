package services

import (
	"time"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func SaveBalance(balance *models.Balance) error {
	return db.DB.Create(balance).Error
}

func GetBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.Balance, error) {
	var balances []models.Balance
	err := db.DB.Where("ecdsa = ? AND eddsa = ?", ecdsaPublicKey, eddsaPublicKey).Find(&balances).Error
	return balances, err
}

func GetLatestBalancesByVaultKeys(ecdsaPublicKey, eddsaPublicKey string) ([]models.Balance, error) {
	var latestBalance models.Balance
	var balances []models.Balance

	err := db.DB.Where("ecdsa = ? AND eddsa = ?", ecdsaPublicKey, eddsaPublicKey).Order("created_at desc").First(&latestBalance).Error
	if err != nil {
		return nil, err
	}

	thresholdTime := latestBalance.CreatedAt.Add(-5 * time.Minute)

	err = db.DB.Where("ecdsa = ? AND eddsa = ? AND created_at >= ?", ecdsaPublicKey, eddsaPublicKey, thresholdTime).Order("created_at desc").Find(&balances).Error
	return balances, err
}

func FetchBalanceOfAddress(address, chain string) (float64, error) {
	return 1.1, nil
}
