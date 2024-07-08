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
