package models

// VaultRequest is the request to add a new vault into registry
type VaultRequest struct {
	Uid                   string `json:"uid" binding:"required"`
	Name                  string `json:"name" binding:"required"`
	PublicKeyECDSA        string `json:"public_key_ecdsa" binding:"required"`
	PublicKeyEDDSA        string `json:"public_key_eddsa" binding:"required"`
	HexChainCode          string `json:"hex_chain_code" binding:"required"`
	ShowNameInLeaderboard bool   `json:"show_name_in_leaderboard"`
}

// VaultRequest is the request to add a new vault into registry
type SharedVaultRequest struct {
	Uid            string `json:"uid" binding:"required"`
	PublicKeyECDSA string `json:"public_key_ecdsa" binding:"required"`
	PublicKeyEDDSA string `json:"public_key_eddsa" binding:"required"`
	HexChainCode   string `json:"hex_chain_code" binding:"required"`
	Theme          string `json:"theme" binding:""`
	Logo           string `json:"logo" binding:""`
}
