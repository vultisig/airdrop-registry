package services

import (
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
	var balances []models.Balance

	subquery := db.DB.Table("balances").
		Select("MAX(id) as id").
		Where("ecdsa = ? AND eddsa = ?", ecdsaPublicKey, eddsaPublicKey).
		Group("chain, token")

	err := db.DB.Where("id IN (?)", subquery).Order("created_at desc").Find(&balances).Error
	return balances, err
}

type ChainTokenPair struct {
	Chain string `json:"chain"`
	Token string `json:"token"`
}

func GetUniqueActiveChainTokenPairs() ([]ChainTokenPair, error) {
	var chainTokenPairs []ChainTokenPair

	subquery := db.DB.Table("balances").
		Select("MAX(id) as id").
		Group("chain, token, ecdsa, eddsa")

	err := db.DB.Table("balances").
		Select("DISTINCT chain, token").
		Where("id IN (?) AND balance > 0", subquery).
		Scan(&chainTokenPairs).Error

	return chainTokenPairs, err
}
