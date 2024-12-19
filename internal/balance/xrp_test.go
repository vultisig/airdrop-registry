package balance

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFetchXRPBalanceOfAddress(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a mock response
		response := map[string]interface{}{
			"result": map[string]interface{}{
				"account_data": map[string]interface{}{
					"Account":           "rhmezeHcxx9sv3A69eafEcAeX3EWBmwFGX",
					"Balance":           "10000000",
					"Flags":             0,
					"LedgerEntryType":   "AccountRoot",
					"OwnerCount":        0,
					"PreviousTxnID":     "E883E3275B66838FF4009DC6C0DF21E263B3A53F7ECCB35CBE38383981F7E0E0",
					"PreviousTxnLgrSeq": 92853100,
					"Sequence":          92803995,
					"index":             "97E1C080C8A7AED3726BAAC7F2DA89E2DF135385823B288A8E13C9F246E81483",
				},
				"account_flags": map[string]interface{}{
					"allowTrustLineClawback":       false,
					"defaultRipple":                false,
					"depositAuth":                  false,
					"disableMasterKey":             false,
					"disallowIncomingCheck":        false,
					"disallowIncomingNFTokenOffer": false,
					"disallowIncomingPayChan":      false,
					"disallowIncomingTrustline":    false,
					"disallowIncomingXRP":          false,
					"globalFreeze":                 false,
					"noFreeze":                     false,
					"passwordSpent":                false,
					"requireAuthorization":         false,
					"requireDestinationTag":        false,
				},
				"ledger_current_index": 92872984,
				"queue_data": map[string]interface{}{
					"txn_count": 0,
				},
				"status":    "success",
				"validated": false,
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create a LiquidityPositionResolver instance
	balanceResolver := &BalanceResolver{
		logger:                logrus.WithField("module", "balance_resolver_test").Logger,
		xrpBalanceBaseAddress: mockServer.URL,
	}
	b, err := balanceResolver.FetchXRPBalanceOfAddress("rhmezeHcxx9sv3A69eafEcAeX3EWBmwFGX")
	assert.NoError(t, err)
	assert.Equal(t, float64(10), b)
}
