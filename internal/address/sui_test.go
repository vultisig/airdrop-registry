package address

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vultisig/airdrop-registry/internal/common"
)

var testECDSAPublicKey = "027e897b35aa9f9fff223b6c826ff42da37e8169fae7be57cbd38be86938a746c6"
var testEdDSAPublicKey = "2dff7cf8446bd3829604bc5c2193ec64c43f67e764de3fd4807df759b91426fe"
var testHexChainCode = "57f3f25c4b034ad80016ef37da5b245bfd6187dc5547696c336ff5a66ed7ee0f"

func TestGetSuiAddress(t *testing.T) {
	tests := []struct {
		name  string
		chain common.Chain
		want  string
	}{
		{
			name:  "Sui",
			chain: common.Sui,
			want:  "0x7a4629f9194d10526e80d76be734535bd5581ef37760d6914052d26066a8ff7b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSuiAddress(testEdDSAPublicKey)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			t.Logf("Got: %s", got)
			assert.Equal(t, got, tt.want)
		})
	}
}
