package models

import "gorm.io/gorm"

// New Type for Vault Share Page Appearance
type VaultShareAppearance struct {
	gorm.Model
	VaultID uint   `json:"vault_id" binding:"required" gorm:"not null;uniqueIndex"`
	Theme   string `json:"theme" binding:"required" gorm:"type:varchar(50)"`
	Logo    string `json:"logo" gorm:"type:text"`
}

func (VaultShareAppearance) TableName() string {
	return "vault_share_appearances"
}
