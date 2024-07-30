package services

import (
	"fmt"
	"strings"

	"github.com/vultisig/airdrop-registry/internal/models"
	"github.com/vultisig/airdrop-registry/pkg/db"
)

func RegisterVault(vault *models.Vault) error {
	if err := db.DB.Create(vault).Error; err != nil {
		return fmt.Errorf("failed to register vault: %w", err)
	}
	return nil
}

func GetVault(ecdsa, eddsa string) (*models.Vault, error) {
	ecdsa = strings.ToLower(ecdsa)
	eddsa = strings.ToLower(eddsa)
	var vault models.Vault
	if err := db.DB.Where("ecdsa = ? AND eddsa = ?", ecdsa, eddsa).First(&vault).Error; err != nil {
		return nil, fmt.Errorf("failed to get vault with ECDSA %s and EDDSA %s: %w", ecdsa, eddsa, err)
	}
	return &vault, nil
}

func UpdateVault(vault *models.Vault) error {
	if err := db.DB.Save(vault).Error; err != nil {
		return fmt.Errorf("failed to update vault: %w", err)
	}
	return nil
}

func DeleteVault(ecdsa, eddsa string) error {
	ecdsa = strings.ToLower(ecdsa)
	eddsa = strings.ToLower(eddsa)
	if err := db.DB.Where("ecdsa = ? AND eddsa = ?", ecdsa, eddsa).Delete(&models.Vault{}).Error; err != nil {
		return fmt.Errorf("failed to delete vault with ECDSA %s and EDDSA %s: %w", ecdsa, eddsa, err)
	}
	return nil
}

func GetAllVaults(page, pageSize int) ([]models.Vault, error) {
	var vaults []models.Vault
	query := db.DB.Order("ecdsa, eddsa")

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&vaults).Error; err != nil {
		return nil, fmt.Errorf("failed to get all vaults: %w", err)
	}
	return vaults, nil
}
