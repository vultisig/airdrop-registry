package balance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFetchTronBalanceOfAddress(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `
		{"data":[
		{"owner_permission":{"keys":[{"address":"TNrTj7SizyxBd4G48cLhZeBvJtZgUaCq2D","weight":1}],
		"threshold":1,"permission_name":"owner"},
		"account_resource":{"energy_window_optimized":true,"latest_consume_time_for_energy":1692016326000,"energy_window_size":28800000},
		"active_permission":[{"operations":"7fff1fc0033ec30f000000000000000000000000000000000000000000000000",
		"keys":[{"address":"TNrTj7SizyxBd4G48cLhZeBvJtZgUaCq2D","weight":1}],
		"threshold":1,"id":2,"type":"Active","permission_name":"active"}],
		"address":"418cc1c4862b6dc2f6550385393bda64fcfdb410d4",
		"create_time":1690283418000,"latest_opration_time":1692016326000,
		"free_asset_net_usageV2":[{"value":0,"key":"1004980"},{"value":0,"key":"1004920"},{"value":0,"key":"1004975"},{"value":0,"key":"1004937"},{"value":0,"key":"1004950"}],"assetV2":[{"value":8888888,"key":"1004920"},{"value":100000,"key":"1004975"},{"value":59860,"key":"1004950"},{"value":17777776,"key":"1004980"},{"value":1777777776,"key":"1004937"}],"frozenV2":[{},{"type":"ENERGY"},{"type":"TRON_POWER"}],
		"balance":26000000,
		"trc20":[{"TDEmUiWoetqL3oGXaxt7F2HqpfWKvHgH7D":"9558000000"},{"TK2y7RAgVhh8WFqN7v4totmLM7ExKc5YtS":"21146259450"},{"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t":"2100000"},{"TTL3B4jgMtBssXzLEiVGzVunWYGBYUBzVx":"3514810000"}],"latest_consume_free_time":1692016326000,"net_window_size":28800000,"net_window_optimized":true}],
		"success":true,"meta":{"at":1738649101361,"page_size":1}}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	// Create a LiquidityPositionResolver instance
	balanceResolver := &BalanceResolver{
		logger:                 logrus.WithField("module", "balance_resolver_test").Logger,
		tronBalanceBaseAddress: mockServer.URL,
	}
	trxBalance, err := balanceResolver.FetchTronBalanceOfAddress("TNrTj7SizyxBd4G48cLhZeBvJtZgUaCq2D", "", 6)
	assert.NoError(t, err)
	assert.Equal(t, float64(26), trxBalance)

	trxBalance, err = balanceResolver.FetchTronBalanceOfAddress("TNrTj7SizyxBd4G48cLhZeBvJtZgUaCq2D", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", 6)
	assert.NoError(t, err)
	assert.Equal(t, float64(2.1), trxBalance)
}
