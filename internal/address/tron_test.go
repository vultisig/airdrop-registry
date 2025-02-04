package address

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/mobile-tss-lib/tss"
)

func TestGetTronAddress(t *testing.T) {
	tests := []struct {
		name  string
		chain common.Chain
		want  string
	}{
		{
			name:  "Tron",
			chain: common.Tron,
			want:  "THFxtPNvc7R9rz4ecC6aTSyPp2WoZnrZh3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.chain.GetDerivePath())
			childPublicKey, err := tss.GetDerivedPubKey(testECDSAPublicKey, testHexChainCode, tt.chain.GetDerivePath(), false)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			got, err := GetTronAddress(childPublicKey)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

// C
