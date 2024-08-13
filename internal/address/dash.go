package address

import (
	"encoding/hex"
	"fmt"

	dashutil "github.com/dashpay/dashd-go/btcutil"
	"github.com/dashpay/dashd-go/chaincfg"
)

func GetDashAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid derived ECDSA public key: %w", err)
	}
	witnessProgram := dashutil.Hash160(pubKeyBytes)
	conv, err := dashutil.NewAddressPubKeyHash(witnessProgram, &chaincfg.MainNetParams)
	if err != nil {
		return "", fmt.Errorf("fail to get public key hash: %w", err)
	}
	return conv.EncodeAddress(), nil
}
