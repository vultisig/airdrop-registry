package address

import (
	"encoding/hex"
	"fmt"

	"github.com/eager7/dogd/chaincfg"
	"github.com/eager7/dogutil"
)

func GetDogeAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid derived ECDSA public key: %w", err)
	}
	witnessProgram := dogutil.Hash160(pubKeyBytes)
	conv, err := dogutil.NewAddressPubKeyHash(witnessProgram, &chaincfg.MainNetParams)
	if err != nil {
		return "", fmt.Errorf("fail to get public key hash: %w", err)
	}
	return conv.EncodeAddress(), nil
}
