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
		result, err := HexToFloat64(tc.hexStr)
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
