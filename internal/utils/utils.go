package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/mr-tron/base58"
)

func IsValidHex(s string) bool {
	// hex 64-66 characters
	re := regexp.MustCompile(`^[0-9a-fA-F]{64,66}$`)
	return re.MatchString(s)
}

func HexToFloat64(hexStr string, decimals int64) (float64, error) {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	value := new(big.Int)
	_, ok := value.SetString(hexStr, 16)
	if !ok {
		return 0, fmt.Errorf("invalid hexadecimal string")
	}

	fValue := new(big.Float).SetInt(value)
	result, _ := new(big.Float).Quo(fValue, new(big.Float).SetInt(math.BigPow(10, decimals))).Float64()
	return result, nil
}

func DecodeBase58ToHex(base58Address string) (string, error) {
	rawBytes, err := base58.Decode(base58Address)
	if err != nil {
		return "", err
	}

	payload := rawBytes[:len(rawBytes)-4]
	hexAddress := fmt.Sprintf("%X", payload)
	return hexAddress, nil
}

func HexToBase58(hexAddress string) (string, error) {
	rawBytes, err := hex.DecodeString(hexAddress)
	if err != nil {
		return "", err
	}

	hash1 := sha256.Sum256(rawBytes)
	hash2 := sha256.Sum256(hash1[:])
	checksum := hash2[:4]

	addressBytes := append(rawBytes, checksum...)
	base58Address := base58.Encode(addressBytes)
	return base58Address, nil
}

func HexToBytes(hexStr string) ([]byte, error) {
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	return hex.DecodeString(hexStr)
}
