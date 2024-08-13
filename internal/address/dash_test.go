package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vultisig/mobile-tss-lib/tss"

	"github.com/vultisig/airdrop-registry/internal/common"
)

func TestGetDashAddress(t *testing.T) {
	tests := []struct {
		name  string
		chain common.Chain
		want  string
	}{
		{
			name:  "Dashcoin",
			chain: common.Dash,
			want:  "XjR1tLaNcAXQUMRVV3w69cM2DMN5aB5gg3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			childPublicKey, err := tss.GetDerivedPubKey(testECDSAPublicKey, testHexChainCode, tt.chain.GetDerivePath(), false)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			got, err := GetDashAddress(childPublicKey)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			t.Logf("Got: %s", got)
			assert.Equal(t, got, tt.want)
		})
	}
}
