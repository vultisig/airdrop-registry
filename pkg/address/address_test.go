package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testECDSAPublicKey = "027e897b35aa9f9fff223b6c826ff42da37e8169fae7be57cbd38be86938a746c6"

// var testEDDSAPublicKey = "2dff7cf8446bd3829604bc5c2193ec64c43f67e764de3fd4807df759b91426fe"
var testHexChainCode = "57f3f25c4b034ad80016ef37da5b245bfd6187dc5547696c336ff5a66ed7ee0f"

func TestGenerateSupportedChainAddresses(t *testing.T) {
	addresses, err := GenerateSupportedChainAddresses(testECDSAPublicKey, testHexChainCode)
	assert.NoError(t, err)

	expectedAddresses := map[string]string{
		"bitcoin":   "bc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7krm4v5f",
		"ethereum":  "0x77435f412e594Fe897fc889734b4FC7665359097",
		"thorchain": "thor1uyhkx5l98awp0q32qqmsx0h440t5cd99q8l3n5",
		"mayachain": "maya1uyhkx5l98awp0q32qqmsx0h440t5cd99qspa9y",
	}

	for chain, expectedAddress := range expectedAddresses {
		assert.Equal(t, expectedAddress, addresses[chain], "Mismatch for %s address", chain)
	}
}

func TestGenerateSpecificAddresses(t *testing.T) {
	testCases := []struct {
		name            string
		chain           string
		derivePath      string
		expectedAddress string
	}{
		{
			name:            "Bitcoin",
			chain:           "bitcoin",
			derivePath:      "m/84'/0'/0'/0/0",
			expectedAddress: "bc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7krm4v5f",
		},
		{
			name:            "Ethereum",
			chain:           "ethereum",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "THORChain",
			chain:           "thorchain",
			derivePath:      "m/44'/931'/0'/0/0",
			expectedAddress: "thor1uyhkx5l98awp0q32qqmsx0h440t5cd99q8l3n5",
		},
		{
			name:            "MayaChain",
			chain:           "mayachain",
			derivePath:      "m/44'/931'/0'/0/0",
			expectedAddress: "maya1uyhkx5l98awp0q32qqmsx0h440t5cd99qspa9y",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			keys, err := GenerateChainKeys(tc.chain, testECDSAPublicKey, testHexChainCode, tc.derivePath, false)
			assert.NoError(t, err)
			assert.Equal(t, tc.chain, keys.ChainName)
			assert.NotEmpty(t, keys.PublicKey)
			assert.Equal(t, tc.expectedAddress, keys.Address)
		})
	}
}
