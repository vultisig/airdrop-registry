package address

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/mobile-tss-lib/tss"
)

func TestGetXRPAddress(t *testing.T) {
	tests := []struct {
		name  string
		chain common.Chain
		want  string
	}{
		{
			name:  "XRP",
			chain: common.XRP,
			want:  "rhmezeHcxx9sv3A69eafEcAeX3EWBmwFGX",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			childPublicKey, err := tss.GetDerivedPubKey(testECDSAPublicKey, testHexChainCode, tt.chain.GetDerivePath(), false)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			got, err := GetXRPAddress(childPublicKey)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			t.Logf("Got: %s", got)
			assert.Equal(t, got, tt.want)
		})
	}
}
