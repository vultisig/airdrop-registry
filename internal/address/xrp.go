package address

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

// Base58 alphabet used by XRP
const alphabet = "rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz"

// AccountID prefix for XRP addresses
var accountIDPrefix = []byte{0x00}

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
	return base58Encode(finalHash), nil
}

func base58Encode(input []byte) string {
	// Count leading zeros
	zeros := 0
	for zeros < len(input) && input[zeros] == 0 {
		zeros++
	}
	result := make([]byte, len(input)*2)
	resultLen := 0
	for _, b := range input {
		carry := int(b)
		for j := 0; j < resultLen || carry != 0; j++ {
			if j > resultLen-1 {
				resultLen++
			}
			carry += int(result[j]) * 256
			result[j] = byte(carry % 58)
			carry /= 58
		}
	}
	encoded := make([]byte, zeros+resultLen)
	for i := 0; i < zeros; i++ {
		encoded[i] = alphabet[0]
	}
	for i := 0; i < resultLen; i++ {
		encoded[zeros+resultLen-1-i] = alphabet[result[i]]
	}
	return string(encoded)
}
