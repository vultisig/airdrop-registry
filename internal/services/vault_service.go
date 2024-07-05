package services

import (
	"strings"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func RegisterVault(vault *models.Vault) error {
	return db.DB.Create(vault).Error
}

func GetVault(eccdsaAddress, eddsaAddress string) (*models.Vault, error) {
	eccdsaAddress = strings.ToLower(eccdsaAddress)
	eddsaAddress = strings.ToLower(eddsaAddress)

	var vault models.Vault
	if err := db.DB.Where("ecdsa = ? AND eddsa = ?", eccdsaAddress, eddsaAddress).First(&vault).Error; err != nil {
		return nil, err
	}
	return &vault, nil
}

func GetAllVaults() ([]models.Vault, error) {
	var vaults []models.Vault
	if err := db.DB.Find(&vaults).Error; err != nil {
		return nil, err
	}
	return vaults, nil
}
