package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vultisig/mobile-tss-lib/tss"

	"github.com/vultisig/airdrop-registry/internal/common"
)

func TestGetBech32Address(t *testing.T) {
	tests := []struct {
		name  string
		chain common.Chain
		hrp   string
		want  string
	}{
		{
			name:  "THORChain",
			chain: common.THORChain,
			hrp:   "thor",
			want:  "thor1uyhkx5l98awp0q32qqmsx0h440t5cd99q8l3n5",
		},
		{
			name:  "MayaChain",
			chain: common.MayaChain,
			hrp:   "maya",
			want:  "maya1uyhkx5l98awp0q32qqmsx0h440t5cd99qspa9y",
		},
		{
			name:  "Cosmos",
			chain: common.GaiaChain,
			hrp:   "cosmos",
			want:  "cosmos13myywet4x5nyhyusp0hq5kyf6fzrlp593u26dx",
		},
		{
			name:  "Kujira",
			chain: common.Kujira,
			hrp:   "kujira",
			want:  "kujira13myywet4x5nyhyusp0hq5kyf6fzrlp59q5gzqv",
		},
		{
			name:  "Terra",
			chain: common.Terra,
			hrp:   "terra",
			want:  "terra15k6d28jvcv5hd989j2g6tk6jus2tk8xzlcph00",
		},
		{
			name:  "TerraClassic",
			chain: common.TerraClassic,
			hrp:   "terra",
			want:  "terra15k6d28jvcv5hd989j2g6tk6jus2tk8xzlcph00",
		},
		{
			name:  "Osmosis",
			chain: common.Osmosis,
			hrp:   "osmo",
			want:  "osmo13myywet4x5nyhyusp0hq5kyf6fzrlp59e8e2m5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			childPublicKey, err := tss.GetDerivedPubKey(testECDSAPublicKey, testHexChainCode, tt.chain.GetDerivePath(), false)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			got, err := GetBech32Address(childPublicKey, tt.hrp)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			assert.Equal(t, got, tt.want)
		})
	}
}
