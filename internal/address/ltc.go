package address

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	ltcchaincfg "github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/ltcutil"
)

func GetLitecoinAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid derived ECDSA public key: %w", err)
	}
	witnessProgram := btcutil.Hash160(pubKeyBytes)
	conv, err := ltcutil.NewAddressWitnessPubKeyHash(witnessProgram, &ltcchaincfg.MainNetParams)
	if err != nil {
		return "", fmt.Errorf("fail to get witness public key hash: %w", err)
	}
	return conv.EncodeAddress(), nil
}
