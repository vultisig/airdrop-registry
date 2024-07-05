package address

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip32"
	"github.com/vultisig/airdrop-registry/pkg/address"
)

var testECDSAPublicKey = "027f1aa43c57848d0f71aa22b055cb9b942132ec73e32e9de5a2419e52cbe026f0"
var testHexChainCode = "32aaa90ca22780d06d38c01e21c07dad5569f5e3b41803db638273466d065978"

func TestGenerateChainKeys(t *testing.T) {
	chainCode, err := hex.DecodeString(testHexChainCode)
	assert.NoError(t, err)

	pubKeyBytes, err := hex.DecodeString(testECDSAPublicKey)
	assert.NoError(t, err)

	masterKey := &bip32.Key{
		Key:       pubKeyBytes,
		ChainCode: chainCode,
		IsPrivate: false,
	}

	supportedChains := []struct {
		name            string
		derivePath      string
		expectedAddress string
	}{
		// These are not actual addresses used by the wallet:
		{name: "bitcoin", derivePath: "m/0/0", expectedAddress: "bc1q88f7eht5rr7jfj8nxu3mnr50z8drndme0qqa5j"},
		{name: "ethereum", derivePath: "m/0/0", expectedAddress: "expected-ethereum-address"},
		{name: "thorchain", derivePath: "m/0/0", expectedAddress: "thor188f7eht5rr7jfj8nxu3mnr50z8drndmela8d9q"},
		{name: "mayachain", derivePath: "m/0/0", expectedAddress: "thor188f7eht5rr7jfj8nxu3mnr50z8drndmela8d9q"},
		// These are
		// {name: "bitcoin", derivePath: "m/84'/0'/0'/0/0", expectedAddress: "expected-bitcoin-address"},
		// {name: "ethereum", derivePath: "m/44'/60'/0'/0/0", expectedAddress: "expected-ethereum-address"},
		// {name: "thorchain", derivePath: "m/44'/931'/0'", expectedAddress: "expected-thorchain-address"},
		// {name: "mayachain", derivePath: "m/44'/931'/0'", expectedAddress: "expected-mayachain-address"},
	}

	for _, chain := range supportedChains {
		keys, err := address.GenerateChainKeys(chain.name, chain.derivePath, masterKey)
		assert.NoError(t, err)
		assert.Equal(t, chain.name, keys.ChainName)
		assert.NotEmpty(t, keys.PublicKey)
		assert.Equal(t, chain.expectedAddress, keys.Address)
	}
}
