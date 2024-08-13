package address

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetBech32Address returns the bech32 address of the given hex public key
func GetBech32Address(hexPublicKey string, hrp string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("fail to decode hex public key,err: %w", err)
	}
	pubKeyHash := btcutil.Hash160(pubKeyBytes)
	return sdk.Bech32ifyAddressBytes(hrp, pubKeyHash)
}
