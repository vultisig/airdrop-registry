package utils

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

func IsValidHex(s string) bool {
	// hex 64-66 characters
	re := regexp.MustCompile(`^[0-9a-fA-F]{64,66}$`)
	return re.MatchString(s)
}

func HexToFloat64(hexStr string) (float64, error) {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	value := new(big.Int)
	_, ok := value.SetString(hexStr, 16)
	if !ok {
		return 0, fmt.Errorf("invalid hexadecimal string")
	}
	fValue := new(big.Float).SetInt(value)
	result, _ := new(big.Float).Quo(fValue, big.NewFloat(1e18)).Float64()
	return result, nil
}
