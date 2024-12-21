package address

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// Base58 alphabet used by XRP
const xrpAlphabet = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"

// cosmos/btcutil/base58 alphabet
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

func GetXRPAddress(hexPublicKey string) (string, error) {
	publicKey, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid hex public key: %v", err)
	}
	sha := sha256.New()
	sha.Write(publicKey)
	hash := sha.Sum(nil)

	ripemd := ripemd160.New()
	ripemd.Write(hash)
	hash = ripemd.Sum(nil)

	versionHash := append([]byte{0}, hash...)

	sha = sha256.New()
	sha.Write(versionHash)
	hash = sha.Sum(nil)

	sha = sha256.New()
	sha.Write(hash)
	hash = sha.Sum(nil)

	checksum := hash[:4]

	finalHash := append(versionHash, checksum...)
	base58Addr := base58.Encode([]byte(finalHash))
	result := ""
	for _, b := range base58Addr {
		index := strings.Index(base58Alphabet, string(b))
		if index == -1 || index >= len(xrpAlphabet) {
			return "", fmt.Errorf("invalid base58 character: %s", string(b))
		}
		result += string(xrpAlphabet[index])
	}
	return result, nil
}
