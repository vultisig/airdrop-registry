package models

type ReferralsAPIResponse struct {
	Total int        `json:"total"`
	Items []Referral `json:"items"`
}

type Referral struct {
	WalletPublicKeyEddsa string `json:"wallet_public_key_eddsa"`
	WalletPublicKeyEcdsa string `json:"wallet_public_key_ecdsa"`
	WalletHexChainCode   string `json:"wallet_hex_chain_code"`
	ParentID             string `json:"parent_id"`
}

type ReferralsSummary struct {
	TotalReferrals int `json:"total_referrals"`
	ValidReferrals int `json:"valid_referrals"`
}
