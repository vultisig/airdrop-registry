package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vultisig/airdrop-registry/internal/common"
)

func TestVaultAddress(t *testing.T) {
	vaultAddress := NewVaultAddress(1)
	vaultAddress.SetAddress(common.Ethereum, "0x123")
	vaultAddress.SetAddress(common.THORChain, "thor1")
	assert.Equal(t, "0x123", vaultAddress.GetAddress(common.Ethereum))
	assert.Equal(t, "thor1", vaultAddress.GetAddress(common.THORChain))
	assert.Equal(t, "0x123", vaultAddress.GetEVMAddress())
	assert.Equal(t, "", vaultAddress.GetAddress(common.BscChain))
	assert.Equal(t, 2, len(vaultAddress.GetAllAddress()))
}
