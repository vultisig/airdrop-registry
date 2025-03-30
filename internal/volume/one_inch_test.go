package volume

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//go:embed one_inch_etherscan.json
var etherscanResponse string

//go:embed one_inch_ethplorer.json
var ethplorerResponse string

func TestOneInchVolume(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("address") == "0xa4a4f610e89488eb4ecc6c63069f241a54485269" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(etherscanResponse))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(ethplorerResponse))
		}
	}))
	defer mockServer.Close()
	oneInch := oneInchVolumeTrack{
		logger:           logrus.WithField("module", "vol_oneInch").Logger,
		etherscanbaseUrl: mockServer.URL,
		ethplorerBaseUrl: mockServer.URL,
	}
	expect := map[string]float64{
		"0x121a38277e0ba795edf8cb6be7935a9773e1ac25": 11.696993778690723,
	}
	res, err := oneInch.processVolume(1715879039, 1715889039, "0xa4a4f610e89488eb4ecc6c63069f241a54485269")
	assert.NoErrorf(t, err, "Failed to get: %v", err)
	assert.Equal(t, expect["0x121a38277e0ba795edf8cb6be7935a9773e1ac25"], res["0x121a38277e0ba795edf8cb6be7935a9773e1ac25"])
}
