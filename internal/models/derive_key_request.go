package models

type DerivePublicKeyRequest struct {
	PublicKeyECDSA string `json:"public_key_ecdsa" binding:"required"`
	HexChainCode   string `json:"hex_chain_code" binding:"required"`
	DerivePath     string `json:"derive_path" binding:"required"`
}
