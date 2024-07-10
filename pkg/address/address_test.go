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
		"bitcoin":      "bc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7krm4v5f",
		"ethereum":     "0x77435f412e594Fe897fc889734b4FC7665359097",
		"thorchain":    "thor1uyhkx5l98awp0q32qqmsx0h440t5cd99q8l3n5",
		"mayachain":    "maya1uyhkx5l98awp0q32qqmsx0h440t5cd99qspa9y",
		"arbitrum":     "0x77435f412e594Fe897fc889734b4FC7665359097",
		"avalanche":    "0x77435f412e594Fe897fc889734b4FC7665359097",
		"bsc":          "0x77435f412e594Fe897fc889734b4FC7665359097",
		"base":         "0x77435f412e594Fe897fc889734b4FC7665359097",
		"bitcoin cash": "bitcoincash:qp9zc8q9wy6kcz5fzx3j094p0dnp8gzfws3j7fzhhm",
		"blast chain":  "0x77435f412e594Fe897fc889734b4FC7665359097",
		"cronoschain":  "0x77435f412e594Fe897fc889734b4FC7665359097",
		"dash":         "Xq45bLN9GCkXYDfVbQShGv2raqJemtwDgr",
		"dogecoin":     "DFhHa4f7zq2y6Z9Nw6uSYm3iHC1RoSpFNz",
		"dydx":         "dydx7e897b35aa9f9fff223b6c826ff42da37e8169fae7be57cbd38be86938a746c6",
		"gaia":         "cosmos1uyhkx5l98awp0q32qqmsx0h440t5cd99g3uxkq",
		"kujira":       "kujira1uyhkx5l98awp0q32qqmsx0h440t5cd99q8l3n5",
		"litecoin":     "ltc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7k2m4v5f",
		"optimism eth": "0x77435f412e594Fe897fc889734b4FC7665359097",
		"polkadot":     "polkadot:7e897b35aa9f9fff223b6c826ff42da37e8169fae7be57cbd38be86938a746c6",
		"polygon":      "0x77435f412e594Fe897fc889734b4FC7665359097",
		"solana":       "3DKvVnBUjJcxc3G18WikUUwdP5R5NWLkuVGVfnDLzW5r",
		"sui":          "0x7e897b35aa9f9fff223b6c826ff42da37e8169fae7be57cbd38be86938a746c6",
		"zksync":       "0x77435f412e594Fe897fc889734b4FC7665359097",
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
		// Cosmos
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
		{
			name:            "Gaia",
			chain:           "gaia",
			derivePath:      "m/44'/118'/0'/0/0",
			expectedAddress: "cosmos1uyhkx5l98awp0q32qqmsx0h440t5cd99g3uxkq",
		},
		{
			name:            "Kujira",
			chain:           "kujira",
			derivePath:      "m/44'/118'/0'/0/0",
			expectedAddress: "kujira1uyhkx5l98awp0q32qqmsx0h440t5cd99q8l3n5",
		},
		// EVM
		{
			name:            "Ethereum",
			chain:           "ethereum",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "Arbitrum",
			chain:           "arbitrum",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "Avalanche",
			chain:           "avalanche",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "BSC",
			chain:           "bsc",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "Base",
			chain:           "base",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "Blast Chain",
			chain:           "blast chain",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "Cronos Chain",
			chain:           "cronoschain",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "Optimism ETH",
			chain:           "optimism eth",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "Polygon",
			chain:           "polygon",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		{
			name:            "zkSync",
			chain:           "zksync",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "0x77435f412e594Fe897fc889734b4FC7665359097",
		},
		// UTXO
		{
			name:            "Bitcoin",
			chain:           "bitcoin",
			derivePath:      "m/84'/0'/0'/0/0",
			expectedAddress: "bc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7krm4v5f",
		},
		{
			name:            "Bitcoin Cash",
			chain:           "bitcoin cash",
			derivePath:      "m/44'/145'/0'/0/0",
			expectedAddress: "bitcoincash:qzsvzzkwt9tjl4lv5c4zwks2nse50gqq6scda6xp00",
		},
		{
			name:            "Dash",
			chain:           "dash",
			derivePath:      "m/44'/5'/0'/0/0",
			expectedAddress: "Xq45bLN9GCkXYDfVbQShGv2raqJemtwDgr",
		},
		{
			name:            "Dogecoin",
			chain:           "dogecoin",
			derivePath:      "m/44'/3'/0'/0/0",
			expectedAddress: "DBiwJDqHyaaNUduVFMidqah5mDajBkmPPH",
		},
		{
			name:            "Litecoin",
			chain:           "litecoin",
			derivePath:      "m/84'/2'/0'/0/0",
			expectedAddress: "ltc1qxv03l5rzukwcqgrkea385lw6v85rngpc249vzr",
		},
		// Other
		{
			name:            "dYdX",
			chain:           "dydx",
			derivePath:      "m/44'/60'/0'/0/0",
			expectedAddress: "dydx13myywet4x5nyhyusp0hq5kyf6fzrlp59c9y7d3",
		},
		{
			name:            "Polkadot",
			chain:           "polkadot",
			derivePath:      "m/44'/354'/0'/0/0",
			expectedAddress: "123K3wPFnMXwm7yr3LizgYTkMhMUwiDiG2rbKWRZbf9PiM2a",
		},
		{
			name:            "Solana",
			chain:           "solana",
			derivePath:      "m/44'/501'/0'/0'",
			expectedAddress: "46ZJUzqDR1dxvX7hFWogsAzyAseAwtb1XNGhtCCNCHW5",
		},
		{
			name:            "Sui",
			chain:           "sui",
			derivePath:      "m/44'/784'/0'/0'/0'",
			expectedAddress: "0x7a4629f9194d10526e80d76be734535bd5581ef37760d6914052d26066a8ff7b",
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
