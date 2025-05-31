package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/mr-tron/base58"
)

func IsValidHex(s string) bool {
	// hex 64-66 characters
	re := regexp.MustCompile(`^[0-9a-fA-F]{64,66}$`)
	return re.MatchString(s)
}
func IsETHAddress(s string) bool {
	return common.IsHexAddress(s)
}

func HexToFloat64(hexStr string, decimals int64) (float64, error) {
	if hexStr == "" {
		return 0, nil
	}
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	value := new(big.Int)
	_, ok := value.SetString(hexStr, 16)
	if !ok {
		return 0, fmt.Errorf("invalid hexadecimal string")
	}

	fValue := new(big.Float).SetInt(value)
	result, _ := new(big.Float).Quo(fValue, new(big.Float).SetInt(ethmath.BigPow(10, decimals))).Float64()
	return result, nil
}

func DecodeBase58ToHex(base58Address string) (string, error) {
	rawBytes, err := base58.Decode(base58Address)
	if err != nil {
		return "", err
	}
	if len(rawBytes) < 4 {
		return "", fmt.Errorf("decoded address is too short")
	}

	payload := rawBytes[:len(rawBytes)-4]
	hexAddress := hex.EncodeToString(payload)
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

func NewJsonRPCRequest(method string, params interface{}, id int) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      id,
	}
}

func EIP55Checksum(address string) (string, error) {
	if !IsETHAddress(address) {
		return "", fmt.Errorf("invalid address")
	}
	addr := common.HexToAddress(address)
	return addr.Hex(), nil
}

// MIN(2,1+(LOG(1+referralCount)/LOG(1+500)))
func GetReferralMultiplier(referralCount int64) float64 {
	// Convert referralCount to float64 for math operations
	rc := float64(referralCount)

	// Compute numerator and denominator
	numerator := math.Log(1 + rc)
	denominator := math.Log(1 + 500)

	// Calculate multiplier
	multiplier := 1 + (numerator / denominator)

	// Apply MIN(2, multiplier)
	if multiplier > 2 {
		multiplier = 2
	}

	return multiplier
}

// =1+0.02*SQRT(swapVolume)
func GetSwapVolumeMultiplier(swapVolume float64) float64 {
	// Calculate the multiplier
	multiplier := 1 + 0.02*math.Sqrt(swapVolume)

	return multiplier
}
