package address

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchutil"
)

func GetBitcoinCashAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid derived ECDSA public key: %w", err)
	}
	witnessProgram := btcutil.Hash160(pubKeyBytes)
	conv, err := bchutil.NewAddressPubKeyHash(witnessProgram, &chaincfg.MainNetParams)
	if err != nil {
		return "", fmt.Errorf("fail to get public key hash: %w", err)
	}
	return conv.EncodeAddress(), nil
}
