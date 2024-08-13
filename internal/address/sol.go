package address

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/btcutil/base58"
)

func GetSolAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid derived ECDSA public key: %w", err)
	}
	return base58.Encode(pubKeyBytes), nil
}
