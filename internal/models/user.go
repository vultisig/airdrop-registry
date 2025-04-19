package models

type ReferralsAPIResponse struct {
	Total int        `json:"total"`
	Items []Referral `json:"items"`
}

type Referral struct {
	Id                   string  `json:"id"`
	UUId                 string  `json:"uuid"`
	FirstName            string  `json:"first_name"`
	Username             string  `json:"username"`
	ReferralsCount       int     `json:"referrals_count"`
	WalletUID            string  `json:"wallet_uid"`
	WalletPublicKeyEddsa *string `json:"wallet_public_key_eddsa"`
	WalletPublicKeyEcdsa *string `json:"wallet_public_key_ecdsa"`
	WalletHexChainCode   *string `json:"wallet_hex_chain_code"`
	ParentID             string  `json:"parent_id"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

type ReferralSummary struct {
	TotalReferrals int `json:"total_referrals"`
	ValidReferrals int `json:"valid_referrals"`
}
