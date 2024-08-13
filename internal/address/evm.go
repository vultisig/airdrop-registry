package address

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func GetEVMAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid derived ECDSA public key: %w", err)
	}
	pubKey, err := crypto.DecompressPubkey(pubKeyBytes)
	if err != nil {
		return "", err
	}
	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}
