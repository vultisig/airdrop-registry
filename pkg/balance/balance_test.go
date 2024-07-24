package balance

import (
	"reflect"
	"strconv"
	"testing"
)

func TestFetchTokensWithBalance(t *testing.T) {
	address := "0xaA11EA95475341c4dDb83aF141B01e52500c23d6"
	tokens, err := FetchTokensWithBalance(address, "ethereum")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	checkBalanceGreaterThan(t, tokens, "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", 1000)
	checkBalanceGreaterThan(t, tokens, "0xdac17f958d2ee523a2206206994597c13d831ec7", 1000)
}

func checkBalanceGreaterThan(t *testing.T, tokens map[string]string, addr string, minBalance float64) {
	balanceStr, ok := tokens[addr]
	if !ok {
		t.Errorf("expected token address %s not found", addr)
		return
	}

	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		t.Errorf("error parsing balance for address %s: %v", addr, err)
		return
	}

	if balance <= minBalance {
		t.Errorf("expected balance for token address %s to be greater than %f, got %f", addr, minBalance, balance)
	}
}

func TestGetTokenInfo(t *testing.T) {
	addresses := []string{"0xc3d688b66703497daa19211eedff47f25384cdc3", "0xd01409314acb3b245cea9500ece3f6fd4d70ea30"}
	tokenInfo, err := GetTokenInfo(addresses, "ethereum")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedTokenAddresses := map[string]TokenInfo{
		"0xc3d688b66703497daa19211eedff47f25384cdc3": {
			Address:  "0xc3d688b66703497daa19211eedff47f25384cdc3",
			Symbol:   "cUSDCv3",
			Decimals: 6,
			Name:     "Compound USDC",
			LogoURI:  "https://tokens.1inch.io/0xc3d688b66703497daa19211eedff47f25384cdc3.png",
			Eip2612:  false,
			Tags:     []string{"tokens"},
		},
		"0xd01409314acb3b245cea9500ece3f6fd4d70ea30": {
			Address:  "0xd01409314acb3b245cea9500ece3f6fd4d70ea30",
			Symbol:   "LTO",
			Decimals: 8,
			Name:     "LTO Network",
			LogoURI:  "https://tokens.1inch.io/0xd01409314acb3b245cea9500ece3f6fd4d70ea30.png",
			Eip2612:  false,
			Tags:     []string{"tokens"},
		},
	}

	for addr, expectedInfo := range expectedTokenAddresses {
		info, ok := tokenInfo[addr]
		if !ok {
			t.Errorf("expected token information for address %s not found", addr)
			continue
		}
		compareTokenInfo(t, expectedInfo, info)
	}
}

func compareTokenInfo(t *testing.T, expected, actual TokenInfo) {
	if expected.Address != actual.Address {
		t.Errorf("expected Address %s, got %s", expected.Address, actual.Address)
	}
	if expected.Symbol != actual.Symbol {
		t.Errorf("expected Symbol %s, got %s", expected.Symbol, actual.Symbol)
	}
	if expected.Decimals != actual.Decimals {
		t.Errorf("expected Decimals %d, got %d", expected.Decimals, actual.Decimals)
	}
	if expected.Name != actual.Name {
		t.Errorf("expected Name %s, got %s", expected.Name, actual.Name)
	}
	if expected.LogoURI != actual.LogoURI {
		t.Errorf("expected LogoURI %s, got %s", expected.LogoURI, actual.LogoURI)
	}
	if expected.Eip2612 != actual.Eip2612 {
		t.Errorf("expected Eip2612 %v, got %v", expected.Eip2612, actual.Eip2612)
	}
	if !reflect.DeepEqual(expected.Tags, actual.Tags) {
		t.Errorf("expected Tags %v, got %v", expected.Tags, actual.Tags)
	}
}
