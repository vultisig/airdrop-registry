package balance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/vultisig/airdrop-registry/internal/common"
)

func TestGetUtxoBalances(t *testing.T) {
	t.Skip()
	b, err := NewBalanceResolver()
	assert.Nil(t, err)
	result, err := b.FetchSolanaBalanceOfAddress("H7FmBYGBi5EmbJaKA88yBgmyGm7eSFdkzCtigwkeaXxb")
	assert.Nil(t, err)
	fmt.Println("Solana balance H7FmBYGBi5EmbJaKA88yBgmyGm7eSFdkzCtigwkeaXxb:", result)

	result, err = b.FetchSuiBalanceOfAddress("0x156e6f6a3f8615008b79dd4871f658ec0da6d70a8540d9dd4d12023b8017e638")
	assert.Nil(t, err)
	fmt.Println("SUI balance for 0x156e6f6a3f8615008b79dd4871f658ec0da6d70a8540d9dd4d12023b8017e638:", result)

	result, err = b.FetchEvmBalanceOfAddress(common.Ethereum, "0x07773707BdA78aC4052f736544928b15dD31c5cc")
	assert.Nil(t, err)
	fmt.Println("ETH balance for 0x07773707BdA78aC4052f736544928b15dD31c5cc:", result)

	result, err = b.fetchERC20TokenBalance(common.Ethereum,
		"0xdac17f958d2ee523a2206206994597c13d831ec7",
		"0x07773707BdA78aC4052f736544928b15dD31c5cc", 6)
	assert.Nil(t, err)
	fmt.Println("USDT balance for 0x07773707BdA78aC4052f736544928b15dD31c5cc:", result)
	balance, balanceUSD, err := b.FetchUtxoBalanceOfAddress("bc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7krm4v5f", common.Bitcoin)
	fmt.Println(balance)
	fmt.Println(balanceUSD)

	balance, err = b.FetchThorchainBalanceOfAddress("thor1tgxm5jw6hrlvslrd6lqpk4jwuu4g29dxytrean")
	assert.Nil(t, err)
	fmt.Println("thor1tgxm5jw6hrlvslrd6lqpk4jwuu4g29dxytrean:", balance)

	balance, err = b.FetchThorchainBalanceOfAddress("thor13amyx54c7z8vfhtd4fhghl30rz2v4t0hdsuk6w")
	assert.Nil(t, err)
	fmt.Println("thor13amyx54c7z8vfhtd4fhghl30rz2v4t0hdsuk6w:", balance)

	balance, err = b.FetchMayachainCacoBalanceOfAddress("maya1h5rlf94hqkvvkyzyhmmgw0hdtw200nqjmaymqc")
	assert.Nil(t, err)
	fmt.Println("maya1h5rlf94hqkvvkyzyhmmgw0hdtw200nqjmaymqc:", balance)

	balance, err = b.FetchMayachainCacoBalanceOfAddress("maya1vzltn37rqccwk95tny657au9j2z072dhg845dr")
	assert.Nil(t, err)
	fmt.Println("maya1vzltn37rqccwk95tny657au9j2z072dhg845dr:", balance)

	balance, err = b.FetchCosmosBalanceOfAddress("cosmos1jl8v454zpnjz76djzdydeq8gwk9364gjked53g")
	assert.Nil(t, err)
	fmt.Println("cosmos1jl8v454zpnjz76djzdydeq8gwk9364gjked53g:", balance)

	balance, err = b.FetchDydxBalanceOfAddress("dydx1jl8v454zpnjz76djzdydeq8gwk9364gjlqrs3l")
	assert.Nil(t, err)
	fmt.Println("dydx1jl8v454zpnjz76djzdydeq8gwk9364gjlqrs3l:", balance)

	balance, err = b.FetchKujiraBalanceOfAddress("kujira153nnvyxz66sj4ywldvy0uexhdnwpfw9fyf4nkz", "ukuji", 6)
	assert.Nil(t, err)
	fmt.Println("kujira153nnvyxz66sj4ywldvy0uexhdnwpfw9fyf4nkz", balance)
}

func TestFetchUtxoBalances(t *testing.T) {
	tests := []struct {
		name         string
		address      string
		chain        common.Chain
		mockResponse UtxoResult
		wantBalance  float64
		wantUSDValue float64
		wantErr      bool
	}{
		{
			name:    "successful bitcoin balance fetch",
			address: "bc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7krm4v5f",
			chain:   common.Bitcoin,
			mockResponse: UtxoResult{
				Data: map[string]struct {
					Address struct {
						Balance    float64 `json:"balance"`
						BalanceUSD float64 `json:"balance_usd"`
					} `json:"address"`
				}{
					"bc1qxpeg8k8xrygj9ae8q6pkzj29sf7w8e7krm4v5f": {
						Address: struct {
							Balance    float64 `json:"balance"`
							BalanceUSD float64 `json:"balance_usd"`
						}{
							Balance:    3934,
							BalanceUSD: 4.28896482,
						},
					},
				},
			},
			wantBalance:  0.00003934,
			wantUSDValue: 4.28896482,
		},
		{
			name:    "successful zcash balance fetch",
			address: "t1M6wQpBni81cypEEMYmrj241TvtyGgLdCu",
			chain:   common.Zcash,
			mockResponse: UtxoResult{
				Data: map[string]struct {
					Address struct {
						Balance    float64 `json:"balance"`
						BalanceUSD float64 `json:"balance_usd"`
					} `json:"address"`
				}{
					"t1M6wQpBni81cypEEMYmrj241TvtyGgLdCu": {
						Address: struct {
							Balance    float64 `json:"balance"`
							BalanceUSD float64 `json:"balance_usd"`
						}{
							Balance:    3238713,
							BalanceUSD: 1.3576684896,
						},
					},
				},
			},
			wantBalance:  0.03238713,
			wantUSDValue: 1.3576684896,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer mockServer.Close()

			resolver := &BalanceResolver{
				logger:           logrus.WithField("module", "balance_resolver_test").Logger,
				vultisigApiProxy: mockServer.URL,
			}

			balance, balanceUSD, err := resolver.FetchUtxoBalanceOfAddress(tt.address, tt.chain)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantBalance, balance)
			assert.Equal(t, tt.wantUSDValue, balanceUSD)
		})
	}
}
