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
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "Valid lowercase address",
			input:          "0x5aeda56215b167893e80b4fe645ba6d5bab767de",
			expectedOutput: "0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE",
			expectError:    false,
		},
		{
			name:           "Valid uppercase address",
			input:          "0X5AEDA56215B167893E80B4FE645BA6D5BAB767DE",
			expectedOutput: "0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE",
			expectError:    false,
		},
		{
			name:           "Address with mixed case",
			input:          "0x5aeDa56215b167893e80B4fE645BA6d5Bab767DE",
			expectedOutput: "0x5AEDA56215b167893e80B4fE645BA6d5Bab767DE",
			expectError:    false,
		},
		{
			name:        "Invalid address - too short",
			input:       "0x5aeda56215b167",
			expectError: true,
		},
		{
			name:        "Invalid address - non-hex characters",
			input:       "0xZZeda56215b167893e80b4fe645ba6d5bab767de",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := EIP55Checksum(tc.input)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for input %s, but got none", tc.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Did not expect error for input %s, but got: %v", tc.input, err)
				return
			}

			if actual != tc.expectedOutput {
				t.Errorf("For input %s, expected output %s, but got %s", tc.input, tc.expectedOutput, actual)
			}
		})
	}
}

func TestGetReferralMultiplier(t *testing.T) {
	testCases := []struct {
		input          int64
		expectedOutput float64
	}{
		{0, 1},
		{1, 1.1114992922647913},
		{10, 1.3857241771164996},
		{500, 2},
	}
	for _, tc := range testCases {
		result := GetReferralMultiplier(tc.input)
		if result != tc.expectedOutput {
			t.Errorf("Expected %f for input %d, but got %f", tc.expectedOutput, tc.input, result)
		}
	}
}

func TestGetSwapVolumeMultiplier(t *testing.T) {
	testCases := []struct {
		input          float64
		expectedOutput float64
	}{
		{0, 1},
		{400, 1.4},
		{900, 1.6},
		{1600, 1.8},
		{2500, 2},
	}
	for _, tc := range testCases {
		result := GetSwapVolumeMultiplier(tc.input)
		if result != tc.expectedOutput {
			t.Errorf("Expected %f for input %f, but got %f", tc.expectedOutput, tc.input, result)
		}
	}
}
