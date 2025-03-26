package tokens

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/vultisig/airdrop-registry/internal/common"
	"github.com/vultisig/airdrop-registry/internal/models"
)

func TestTrcAutoDiscovery(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBody   trcAccountResponse
		expectedCoins  []models.Coin
		expectedError  bool
	}{
		{
			name:           "successful response with TRX balance",
			responseStatus: http.StatusOK,
			responseBody: trcAccountResponse{
				Success: true,
				Data: []trcAccount{
					{
						Balance: 1000000,
						Trc20:   []map[string]string{},
					},
				},
			},
			expectedCoins: []models.Coin{
				{
					ContractAddress: "TRX",
					Decimals:        1000000,
				},
			},
			expectedError: false,
		},
		{
			name:           "unsuccessful response",
			responseStatus: http.StatusOK,
			responseBody: trcAccountResponse{
				Success: false,
			},
			expectedCoins: nil,
			expectedError: false,
		},
		{
			name:           "invalid status code",
			responseStatus: http.StatusBadRequest,
			responseBody:   trcAccountResponse{},
			expectedCoins:  nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Create service with test server URL
			url := "https://api.trongrid.io/v1"
			svc := &trcDiscoveryService{
				logger:       logrus.New(),
				tronBaseURL:  url,
				cmcIDService: NewCMCIDService(),
			}

			// Call method
			coins, err := svc.discover("TS98H6jSx6uv1gG1vx6CJZMeYGkMZXgQ7K", common.Tron)

			// Check results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedCoins == nil {
					assert.Nil(t, coins)
				} else {
					assert.Equal(t, len(tt.expectedCoins), len(coins))
					for i := range tt.expectedCoins {
						assert.Equal(t, tt.expectedCoins[i].ContractAddress, coins[i].ContractAddress)
						assert.Equal(t, tt.expectedCoins[i].Decimals, coins[i].Decimals)
					}
				}
			}
		})
	}
}

func TestFetchTokenData(t *testing.T) {
	decimalResponse := `{"result":{"result":true},"energy_used":2207,"constant_result":["0000000000000000000000000000000000000000000000000000000000000006"],"energy_penalty":1699,"transaction":{"ret":[{}],"visible":false,"txID":"a409443e88c91a982fbd76869ffb556c2af0e3da3570fd7df040663852068c66","raw_data":{"contract":[{"parameter":{"value":{"data":"313ce567","owner_address":"41977c20977f412c2a1aa4ef3d49fee5ec4c31cdfb","contract_address":"41a614f803b6fd780986a42c78ec9c7f77e6ded13c"},"type_url":"type.googleapis.com/protocol.TriggerSmartContract"},"type":"TriggerSmartContract"}],"ref_block_bytes":"8d31","ref_block_hash":"dcfcd108763da59f","expiration":1742901750000,"timestamp":1742901692497},"raw_data_hex":"0a028d312208dcfcd108763da59f40f0d1a8e8dc325a6d081f12690a31747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e54726967676572536d617274436f6e747261637412340a1541977c20977f412c2a1aa4ef3d49fee5ec4c31cdfb121541a614f803b6fd780986a42c78ec9c7f77e6ded13c2204313ce56770d190a5e8dc32"}}`
	symbolResponse := `{"result":{"result":true},"energy_used":5922,"constant_result":["000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000045553445400000000000000000000000000000000000000000000000000000000"],"energy_penalty":4558,"transaction":{"ret":[{}],"visible":false,"txID":"47cde9c904c7b0692a19b72a2adfcbcd86b4e9e8c22cd1c19f48bff5dc365fa1","raw_data":{"contract":[{"parameter":{"value":{"data":"95d89b41","owner_address":"41b1623aaeeac63781b6b2b2325ee6d434b6af4b24","contract_address":"41a614f803b6fd780986a42c78ec9c7f77e6ded13c"},"type_url":"type.googleapis.com/protocol.TriggerSmartContract"},"type":"TriggerSmartContract"}],"ref_block_bytes":"8d43","ref_block_hash":"27c4b850919dfc59","expiration":1742901804000,"timestamp":1742901745991},"raw_data_hex":"0a028d43220827c4b850919dfc5940e0f7abe8dc325a6d081f12690a31747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e54726967676572536d617274436f6e747261637412340a1541b1623aaeeac63781b6b2b2325ee6d434b6af4b24121541a614f803b6fd780986a42c78ec9c7f77e6ded13c220495d89b4170c7b2a8e8dc32"}}`

	tests := []struct {
		name           string
		responseStatus int
		responseBody   string
		expectedSymbol string
		expectedError  bool
	}{
		{
			name:           "successful response with TRX symbol",
			responseStatus: http.StatusOK,
			responseBody:   symbolResponse,
			expectedSymbol: "USDT",
			expectedError:  false,
		},
		{
			name:           "successful response with TRX symbol",
			responseStatus: http.StatusOK,
			responseBody:   decimalResponse,
			expectedSymbol: "USDT",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseStatus)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Create service with test server URL
			url := "https://api.trongrid.io"
			td := &trcDiscoveryService{
				logger:       logrus.New(),
				tronBaseURL:  url,
				cmcIDService: NewCMCIDService(),
			}

			// Call method
			address := "TS98H6jSx6uv1gG1vx6CJZMeYGkMZXgQ7K"
			contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
			selector := "symbol()"
			symbol, err := td.fetchTokenData(address, contract, selector)

			// Check results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSymbol, symbol)
			}
		})
	}

}
