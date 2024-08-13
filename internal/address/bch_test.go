package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vultisig/mobile-tss-lib/tss"

	"github.com/vultisig/airdrop-registry/internal/common"
)

func TestGetBitcoinCashAddress(t *testing.T) {
	tests := []struct {
		name  string
		chain common.Chain
		want  string
	}{
		{
			name:  "BitcoinCash",
			chain: common.BitcoinCash,
			want:  "qzsvzzkwt9tjl4lv5c4zwks2nse50gqq6scda6xp00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			childPublicKey, err := tss.GetDerivedPubKey(testECDSAPublicKey, testHexChainCode, tt.chain.GetDerivePath(), false)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			got, err := GetBitcoinCashAddress(childPublicKey)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
