package models

// VaultResponse to client side(front-end web)
type VaultResponse struct {
	UId                   string        `json:"uid"`
	Name                  string        `json:"name"`
	Alias                 string        `json:"alias"`
	PublicKeyECDSA        string        `json:"public_key_ecdsa"`
	PublicKeyEDDSA        string        `json:"public_key_eddsa"`
	TotalPoints           float64       `json:"total_points"`
	JoinAirdrop           bool          `json:"join_airdrop"`
	Rank                  int64         `json:"rank"`
	SwapVolumeRank        int64         `json:"swap_volume_rank"`
	Balance               int64         `json:"balance"`
	LPValue               int64         `json:"lp_value"`
	NFTValue              int64         `json:"nft_value"`
	Coins                 []ChainCoins  `json:"chains"`
	RegisteredAt          int64         `json:"registered_at"`
	AvatarURL             string        `json:"avatar_url"`
	ShowNameInLeaderboard bool          `json:"show_name_in_leaderboard"`
	SwapVolume            float64       `json:"swap_volume"`
	ReferralCode          string        `json:"referral_code"`
	ReferralCount         int64         `json:"referral_count"`
	SeasonActivities      []SeasonStats `json:"season_stats"` // Needed to highlight user in the leaderboard of each season
}

type SeasonStats struct {
	SeasonID uint    `json:"season_id"`
	Rank     int64   `json:"rank"`
	Points   float64 `json:"points"`
}

type VaultsResponse struct {
	Vaults          []VaultResponse `json:"vaults"`
	TotalVaultCount int64           `json:"total_vault_count"`
	TotalBalance    int64           `json:"total_balance"`
	TotalLP         int64           `json:"total_lp"`
	TotalNFT        int64           `json:"total_nft"`
	TotalSwapVolume float64         `json:"total_swap_volume"`
}
