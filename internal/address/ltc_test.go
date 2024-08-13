package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vultisig/mobile-tss-lib/tss"

	"github.com/vultisig/airdrop-registry/internal/common"
)

func TestGetLitecoinAddress(t *testing.T) {
	tests := []struct {
		name  string
		chain common.Chain
		want  string
	}{
		{
			name:  "Litecoin",
			chain: common.Litecoin,
			want:  "ltc1qxv03l5rzukwcqgrkea385lw6v85rngpc249vzr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			childPublicKey, err := tss.GetDerivedPubKey(testECDSAPublicKey, testHexChainCode, tt.chain.GetDerivePath(), false)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			got, err := GetLitecoinAddress(childPublicKey)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
