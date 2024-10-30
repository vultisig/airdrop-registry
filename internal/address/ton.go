package address

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

var (
	// This is the hex-encoded v3R1 wallet code
	walletV3Code = "B5EE9C72410101010044000084FF0020DDA4F260810200D71820D70B1FED44D0D31FD3FFD15112BAF2A122F901541044F910F2A2F80001D31F3120D74A96D307D402FB00DED1A4C8CB1FCBFFC9ED54"

	// Cell structure for state init
	walletV3Data = "00000000" // Initial data prefix for v3R1
)

func GetTonAddress(hexPublicKey string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(hexPublicKey)
	if err != nil {
		return "", fmt.Errorf("invalid public key: %w", err)
	}
	fmt.Printf("Public Key Bytes: %x\n", pubKeyBytes)

	// Create 36-byte sequence for address
	addressBytes := make([]byte, 36)
	addressBytes[0] = 0x11 // bounceable flag
	addressBytes[1] = 0x00 // workchain_id (0 for basechain)

	// Create state init
	stateInit, err := createWalletStateInit(pubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create state init: %w", err)
	}
	fmt.Printf("State Init: %x\n", stateInit)

	// Hash the state init to get account_id
	accountID := sha256.Sum256(stateInit)
	fmt.Printf("Account ID: %x\n", accountID[:])
	copy(addressBytes[2:34], accountID[:])

	// Add CRC16
	crc := CalculateCRC16(addressBytes[:34])
	fmt.Printf("CRC: %x\n", crc)
	binary.BigEndian.PutUint16(addressBytes[34:], crc)

	fmt.Printf("Final Address Bytes: %x\n", addressBytes)
	result := base64.URLEncoding.EncodeToString(addressBytes)
	fmt.Printf("Final Base64: %s\n", result)

	return result, nil
}

func createWalletStateInit(pubKey []byte) ([]byte, error) {
	// Decode wallet code
	code, err := hex.DecodeString(walletV3Code)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet code: %w", err)
	}

	// Create data cell with public key
	data := make([]byte, 0, len(walletV3Data)/2+len(pubKey))
	dataPrefix, err := hex.DecodeString(walletV3Data)
	if err != nil {
		return nil, fmt.Errorf("invalid data prefix: %w", err)
	}
	data = append(data, dataPrefix...)
	data = append(data, pubKey...)

	// Create state init structure according to TL-B format
	stateInit := make([]byte, 0)

	// Split depth and special flags
	stateInit = append(stateInit, 0x00, 0x00) // No split depth, no special

	// Code cell
	stateInit = append(stateInit, 0x02) // Has code reference
	stateInit = append(stateInit, code...)

	// Data cell
	stateInit = append(stateInit, 0x02) // Has data reference
	stateInit = append(stateInit, data...)

	// Library
	stateInit = append(stateInit, 0x00) // No library

	return stateInit, nil
}

func CalculateCRC16(data []byte) uint16 {
	crc := uint16(0xFFFF)
	for _, b := range data {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc = crc << 1
			}
		}
	}
	return crc
}
