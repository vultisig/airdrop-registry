package address

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vultisig/airdrop-registry/internal/common"
)

func TestGetTonAddress(t *testing.T) {
	// Replace with your public key in hex format
	hexPublicKey := "5a6f496e61121e8679585e81297bd68c01e7946abbfb3eb263753f1d41390fe8"
	walletAddress, err := GetTonAddress(hexPublicKey)
	if err != nil {
		t.Errorf("Failed to create the address: %v", err)
		t.FailNow()
	}

	tests := []struct {
		name  string
		chain common.Chain
		want  string
	}{
		{
			name:  "TON",
			chain: common.Ton,
			want:  "UQA_fNiw1Jrk-TGK2Xknb5_rPTzZGhWPVKcR8ORbNcyTKXEF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, walletAddress)
		})
	}
}
