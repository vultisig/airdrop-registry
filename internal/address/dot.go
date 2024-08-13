package address

import (
	"encoding/hex"
	"fmt"

	"github.com/vultisig/airdrop-registry/pkg/utils"
)

func GetDotAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid derived ECDSA public key: %w", err)
	}
	return utils.SS58Encode(pubKeyBytes, 0), nil
}
