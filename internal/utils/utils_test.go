package utils

import (
	"testing"
)

func TestIsValidHex(t *testing.T) {
	validHexes := []string{
		"027e897b35aa9f9fff223b6c826ff42da37e8169fae7be57cbd38be86938a746c6",
		"57f3f25c4b034ad80016ef37da5b245bfd6187dc5547696c336ff5a66ed7ee0f",
	}

	invalidHexes := []string{
		"invalidhexstring",
		"12345",                // too short
		"ghijklmnopqrstuvwxyz", // contains non-hex characters
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefg", // too long
	}

	for _, hex := range validHexes {
		if !IsValidHex(hex) {
			t.Errorf("Expected %s to be valid", hex)
		}
	}

	for _, hex := range invalidHexes {
		if IsValidHex(hex) {
			t.Errorf("Expected %s to be invalid", hex)
		}
	}
}

func TestHexToFloat64(t *testing.T) {
	testCases := []struct {
		hexStr   string
		expected float64
		hasError bool
	}{
		{"0x027e897b35aa9f9fff223b6c826ff42da37e8169fae7be57cbd38be86938a746c6", 2.888185058056452e+59, false},
		{"57f3f25c4b034ad80016ef37da5b245bfd6187dc5547696c336ff5a66ed7ee0f", 3.978223437431612e+58, false},
		{"invalidhexstring", 0, true},
	}

	for _, tc := range testCases {
		result, err := HexToFloat64(tc.hexStr, 18)
		if tc.hasError {
			if err == nil {
				t.Errorf("Expected error for hex string %s, but got none", tc.hexStr)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect error for hex string %s, but got %v", tc.hexStr, err)
			}
			if result != tc.expected {
				t.Errorf("Expected %v for hex string %s, but got %v", tc.expected, tc.hexStr, result)
			}
		}
	}
}

func TestEIP55(t *testing.T) {
	contracts := make([]string, 0)
	contracts = append(contracts, "0x0b3e328455c4059eeb9e3f84b5543f74e24e7e1b",
		"0x5947bb275c521040051d82396192181b413227a3",
		"0x23ee2343b892b1bb63503a4fabc840e0e2c6810f",
		"tr7nhqjekqxgtci8q8zy4pl8otszgjlj6t",
		"0xfd086bc7cd5c481dcc9c85ebe478a1c0b69fcbb9",
		"es9vmfrzacermjfrf4h2fyd4kconky11mcce8benwnyb",
	)

	for _, contract := range contracts {
		_, err := EIP55Checksum(contract)
		if err != nil {
			t.Errorf("Failed to get EIP55 checksum for %s: %v", contract, err)
		}
	}
}
