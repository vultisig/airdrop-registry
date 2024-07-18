package balance

import (
	"fmt"
	"testing"
)

func TestBitcoinBalance(t *testing.T) {
	bitcoinAddress := "bc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7krm4v5f"
	expectedBalance := 0.0002414

	balance, err := FetchBitcoinBalanceOfAddress(bitcoinAddress)
	if err != nil {
		t.Errorf("Error fetching balance: %v", err)
	}

	if fmt.Sprintf("%.4f", balance) != fmt.Sprintf("%.4f", expectedBalance) {
		t.Errorf("Expected balance %f, got %f", expectedBalance, balance)
	}
}

func TestEthereumBalance(t *testing.T) {
	ethereumAddress := "0x77435f412e594Fe897fc889734b4FC7665359097"
	expectedBalance := 0.054832

	balance, err := FetchEvmBalanceOfAddress("ethereum", ethereumAddress)
	if err != nil {
		t.Errorf("Error fetching balance: %v", err)
	}

	if fmt.Sprintf("%.6f", balance) != fmt.Sprintf("%.6f", expectedBalance) {
		t.Errorf("Expected balance %f, got %f", expectedBalance, balance)
	}
}

func TestTHORChainBalance(t *testing.T) {
	thorchainAddress := "thor1uyhkx5l98awp0q32qqmsx0h440t5cd99q8l3n5"
	expectedBalance := 34.212954

	balance, err := FetchThorchainBalanceOfAddress(thorchainAddress)
	if err != nil {
		t.Errorf("Error fetching balance: %v", err)
	}

	if fmt.Sprintf("%.3f", balance) != fmt.Sprintf("%.3f", expectedBalance) {
		t.Errorf("Expected balance %f, got %f", expectedBalance, balance)
	}
}

func TestMayaChainBalance(t *testing.T) {
	mayaChainAddress := "maya1uyhkx5l98awp0q32qqmsx0h440t5cd99qspa9y"
	expectedBalance := 0.0

	balance, err := FetchMayachainBalanceOfAddress(mayaChainAddress)
	if err != nil {
		t.Errorf("Error fetching balance: %v", err)
	}

	if fmt.Sprintf("%.3f", balance) != fmt.Sprintf("%.3f", expectedBalance) {
		t.Errorf("Expected balance %f, got %f", expectedBalance, balance)
	}
}

func TestPolkadotBalance(t *testing.T) {
	polkadotAddress := "16fq6FSxb8s5Ah2m2wi7mEnemvG7hwithfMqXx6N2FsTumnL"
	expectedBalance := 218.706

	balance, err := FetchPolkadotBalanceOfAddress(polkadotAddress)
	if err != nil {
		t.Errorf("Error fetching balance: %v", err)
	}

	if fmt.Sprintf("%.6f", balance) != fmt.Sprintf("%.6f", expectedBalance) {
		t.Errorf("Expected balance %f, got %f", expectedBalance, balance)
	}
}

// func TestSuiBalance(t *testing.T) {
// 	suiAddress := "0x3a6fD4d2F4d3Cd7f1aBb3E2f1aFbdff4D3D3cd7f"
// 	expectedBalance := 0.987654

// 	balance, err := FetchSuiBalanceOfAddress(suiAddress)
// 	if err != nil {
// 		t.Errorf("Error fetching balance: %v", err)
// 	}

// 	if fmt.Sprintf("%.6f", balance) != fmt.Sprintf("%.6f", expectedBalance) {
// 		t.Errorf("Expected balance %f, got %f", expectedBalance, balance)
// 	}
// }

func TestSolanaBalance(t *testing.T) {
	solanaAddress := "GYRsheZ78JMfMNETuAZNrs6L1U3GsHP5crzzLPeETDYm"
	expectedBalance := 1.5478

	balance, err := FetchSolanaBalanceOfAddress(solanaAddress)
	if err != nil {
		t.Errorf("Error fetching balance: %v", err)
	}

	if fmt.Sprintf("%.2f", balance) != fmt.Sprintf("%.2f", expectedBalance) {
		t.Errorf("Expected balance %f, got %f", expectedBalance, balance)
	}
}
