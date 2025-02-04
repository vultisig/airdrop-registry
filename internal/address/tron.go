package address

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gcash/bchutil/base58"
	"github.com/vultisig/airdrop-registry/internal/utils"
)

func GetTronAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid hex public key: %w", err)
	}
	pubKey, err := crypto.DecompressPubkey(pubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("error decompressing public key: %w", err)
	}
	address := crypto.PubkeyToAddress(*pubKey).Hex()
	address = "41" + address[2:]
	addb, _ := hex.DecodeString(address)
	hash1 := utils.SHA256(utils.SHA256(addb))
	secret := hash1[:4]
	addb = append(addb, secret...)
	return base58.Encode(addb), nil
}
