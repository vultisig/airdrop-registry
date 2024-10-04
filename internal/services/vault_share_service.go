package services

import (
	"fmt"

	"github.com/vultisig/airdrop-registry/internal/models"
)

func (s *Storage) UpdateTheme(appearance models.VaultShareAppearance) error {
	// insert or update theme
	qry := `INSERT INTO vault_share_appearances (vault_id, theme, logo) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE theme = ?, logo = ?`
	if err := s.db.Exec(qry, appearance.VaultID, appearance.Theme, appearance.Logo, appearance.Theme, appearance.Logo).Error; err != nil {
		return fmt.Errorf("failed to update theme: %w", err)
	}
	return nil
}

func (s *Storage) GetTheme(vaultID uint) models.VaultShareAppearance {
	var appearance models.VaultShareAppearance
	s.db.Where("vault_id = ?", vaultID).First(&appearance)
	return appearance
}
