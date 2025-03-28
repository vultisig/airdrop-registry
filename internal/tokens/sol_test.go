package tokens

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

//go:embed sol_token_response.json
var sol_token_response string

func TestSolDiscoveryService_processTokenAccounts(t *testing.T) {
	// Setup
	service := &solanaDiscoveryService{
		logger:         logrus.New(),
		cmcService:     NewCMCService(),
		oneInchService: NewOneinchService(),
	}

	// Parse the embedded JSON response
	var response jsonRpcResponse
	err := json.Unmarshal([]byte(sol_token_response), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal embedded JSON: %v", err)
	}

	testAddr := "CG4V2eoUXnwJSDsmr1fNdbR9r63XHLKD9gA2xpCRdRby"
	results, err := service.processTokenAccounts(&response, testAddr, common.Solana)

	// Assertions
	if err != nil {
		t.Errorf("processTokenAccounts() error = %v", err)
		return
	}

	if len(results) == 0 {
		t.Error("processTokenAccounts() returned empty results")
		return
	}

	expected := models.CoinBase{
		Address:         testAddr,
		Balance:         "24389303",
		Chain:           common.Solana,
		IsNative:        false,
		ContractAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	}

	if !reflect.DeepEqual(results[0].Address, expected.Address) ||
		!reflect.DeepEqual(results[0].Balance, expected.Balance) ||
		!reflect.DeepEqual(results[0].Chain, expected.Chain) ||
		!reflect.DeepEqual(results[0].IsNative, expected.IsNative) ||
		!reflect.DeepEqual(results[0].ContractAddress, expected.ContractAddress) {
		t.Errorf("processTokenAccounts() = %v, want %v", results[0], expected)
	}
}
