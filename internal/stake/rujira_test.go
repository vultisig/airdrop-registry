package stake

import (
	_ "embed"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRujiraAutoStakeBalance(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{
			"balances": []map[string]any{
				{
					"denom":  "bsc-bnb",
					"amount": "1",
				},
				{
					"denom":  "bsc-usdc-0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d",
					"amount": "1",
				},
				{
					"denom":  "rune",
					"amount": "23",
				},
				{
					"denom":  "x/ruji",
					"amount": "1",
				},
				{
					"denom":  "x/staking-x/ruji",
					"amount": "225687",
				}}}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}))

	s := RujiraStakeResolver{
		thornodeBaseAddress: mockServer.URL,
		chainDecimal:        8,
		rujiPrice:           0.862,
	}
	expectedResult := 225687 * math.Pow10(-s.chainDecimal) * 0.862
	stakeBalance, err := s.GetRujiraAutoCompoundStake("thor15dmp7pnhmjslnshh6zszkq2xwmuamyetzn7mn8")
	assert.NoErrorf(t, err, "Failed to get: %v", err)
	assert.Equal(t, expectedResult, stakeBalance)
}

func TestRujiStakeBalance(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]any{"data": map[string]any{"addr": "thor15dmp7pnhmjslnshh6zszkq2xwmuamyetzn7mn8", "bonded": "2456321", "pending_revenue": "0"}}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	s := RujiraStakeResolver{
		thornodeBaseAddress: mockServer.URL,
		chainDecimal:        8,
		rujiPrice:           0.862,
	}
	expectedResult := 2456321 * math.Pow10(-s.chainDecimal) * 0.862
	stakeBalance, err := s.GetRujiraSimpleStake("thor15dmp7pnhmjslnshh6zszkq2xwmuamyetzn7mn8")
	assert.NoErrorf(t, err, "Failed to get: %v", err)
	assert.Equal(t, float64(expectedResult), stakeBalance)
}
