package models

import "github.com/vultisig/airdrop-registry/internal/common"

type VaultAddress struct {
	vaultID uint
	//map of chain name
	address map[common.Chain]string
}

func NewVaultAddress(vaultID uint) VaultAddress {
	return VaultAddress{
		vaultID: vaultID,
		address: make(map[common.Chain]string),
	}
}
func (v *VaultAddress) GetVaultID() uint {
	return v.vaultID
}

func (v *VaultAddress) GetAddress(chain common.Chain) string {
	return v.address[chain]
}
func (v *VaultAddress) GetEVMAddress() string {
	for _, chain := range common.EVMChains {
		if _, ok := v.address[chain]; ok {
			return v.address[chain]
		}
	}
	return ""
}

func (v *VaultAddress) SetAddress(chain common.Chain, address string) {
	v.address[chain] = address
}
func (v *VaultAddress) GetAllAddress() []string {
	var addresses []string
	for _, address := range v.address {
		addresses = append(addresses, address)
	}
	return addresses
}
