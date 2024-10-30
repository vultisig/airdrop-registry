package address

import (
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/pbkdf2"

	"github.com/vultisig/airdrop-registry/internal/common"
)

var mnemonic = "steel address phone tobacco harsh powder denial differ mix jealous kind immune mobile easily stairs ivory original exercise attitude young luggage exotic fresh cost"

func TestGetTonAddress(t *testing.T) {
	seed := MnemonicToSeed(mnemonic)
	pubKey, _ := SeedToKeypair(seed)

	t.Log("Seed:", seed)
	t.Log("Public Key:", pubKey)

	// Convert public key to hex string
	pubKeyHex := hex.EncodeToString(pubKey)
	t.Log("Derived Public Key: ", pubKeyHex)

	tests := []struct {
		name  string
		chain common.Chain
		want  string
	}{
		{
			name:  "TON",
			chain: common.Ton,
			want:  "EQB8V7T2w7BJv1jWC9wVunKJHtSv7Efis_gGEDig1E6-l=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTonAddress(pubKeyHex)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			t.Logf("Got: %s", got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func MnemonicToSeed(mnemonic string) []byte {
	// PBKDF2 parameters for TON
	password := []byte(mnemonic)
	salt := []byte("TON default seed")
	iterations := 100000
	keyLen := 64

	// Generate seed using PBKDF2-HMAC-SHA512
	return pbkdf2.Key(password, salt, iterations, keyLen, sha512.New)
}

func SeedToKeypair(seed []byte) (publicKey, privateKey []byte) {
	// TON uses the first 32 bytes of the seed
	privateKey = ed25519.NewKeyFromSeed(seed[:32])
	publicKey = privateKey[32:] // Ed25519 public key is last 32 bytes
	return publicKey, privateKey
}
