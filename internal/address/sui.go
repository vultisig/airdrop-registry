package address

import (
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

// GetSuiAddress generates a Sui address from a hex-encoded public key
func GetSuiAddress(hexPublicKey string) (string, error) {
	// Decode the hex-encoded public key
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid public key: %w", err)
	}
	toHash := make([]byte, 0, len(pubKeyBytes)+1)
	toHash = append(toHash, 0x00)
	toHash = append(toHash, pubKeyBytes...)
	hashed := blake2b.Sum256(toHash)
	// Convert the hashed public key to a Sui address
	address := hex.EncodeToString(hashed[:])

	return "0x" + address, nil
}
