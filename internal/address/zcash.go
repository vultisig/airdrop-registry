package address

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/vultisig/airdrop-registry/internal/utils"
)

func GetZcashAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid hex public key: %w", err)
	}

	if len(pubKeyBytes) != 33 && len(pubKeyBytes) != 65 {
		return "", fmt.Errorf("invalid public key length: %d bytes", len(pubKeyBytes))
	}

	pubKeyHash := utils.Hash160(pubKeyBytes)

	// Zcash t1 address version prefix: 0x1CB8
	version := []byte{0x1C, 0xB8}
	payload := append(version, pubKeyHash...)

	// Base58Check encoding: payload || checksum
	checksum := sha256.Sum256(payload)
	checksum = sha256.Sum256(checksum[:])
	address := base58.Encode(append(payload, checksum[:4]...))

	return address, nil
}
