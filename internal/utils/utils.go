package utils

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common/math"
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
