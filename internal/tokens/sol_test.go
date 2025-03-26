package tokens

import (
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func TestSolDiscoveryService_processTokenAccounts(t *testing.T) {
	// Setup
	service := &solDiscoveryService{
		logger:         logrus.New(),
		cmcIDService:   NewCMCIDService(),
		oneinchService: NewOneinchService(),
	}

	// Sample response with token accounts
	response := &jsonRpcResponse{
		Jsonrpc: "2.0",
		Result: result{
			Value: []accountEntry{
				{
					Account: account{
						Data: struct {
							Parsed struct {
								Info struct {
									IsNative    bool        `json:"isNative"`
									Mint        string      `json:"mint"`
									Owner       string      `json:"owner"`
									State       string      `json:"state"`
									TokenAmount tokenAmount `json:"tokenAmount"`
								} `json:"info"`
								Type string `json:"type"`
							} `json:"parsed"`
							Program string `json:"program"`
						}{
							Parsed: struct {
								Info struct {
									IsNative    bool        `json:"isNative"`
									Mint        string      `json:"mint"`
									Owner       string      `json:"owner"`
									State       string      `json:"state"`
									TokenAmount tokenAmount `json:"tokenAmount"`
								} `json:"info"`
								Type string `json:"type"`
							}{
								Info: struct {
									IsNative    bool        `json:"isNative"`
									Mint        string      `json:"mint"`
									Owner       string      `json:"owner"`
									State       string      `json:"state"`
									TokenAmount tokenAmount `json:"tokenAmount"`
								}{
									IsNative: false,
									Mint:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
									Owner:    "CG4V2eoUXnwJSDsmr1fNdbR9r63XHLKD9gA2xpCRdRby",
									State:    "initialized",
									TokenAmount: tokenAmount{
										Amount:         "24389303",
										Decimals:       6,
										UIAmount:       24.389303,
										UIAmountString: "24.389303",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	testAddr := "CG4V2eoUXnwJSDsmr1fNdbR9r63XHLKD9gA2xpCRdRby"
	results, err := service.processTokenAccounts(response, testAddr, common.Solana)

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
