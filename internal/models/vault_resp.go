package models

// VaultResponse to client side(front-end web)
type VaultResponse struct {
	UId            string       `json:"uid"`
	Name           string       `json:"name"`
	Alias          string       `json:"alias"`
	PublicKeyECDSA string       `json:"public_key_ecdsa"`
	PublicKeyEDDSA string       `json:"public_key_eddsa"`
	TotalPoints    float64      `json:"total_points"`
	JoinAirdrop    bool         `json:"join_airdrop"`
	Coins          []ChainCoins `json:"chains"`
}
