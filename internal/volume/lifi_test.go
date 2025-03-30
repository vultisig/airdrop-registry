package volume

import (
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//go:embed lifi_test_response.json
var lisfiResponse string

func TestLifiVolume(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(lisfiResponse))
	}))
	defer mockServer.Close()
	li := lifiVolumeTracker{
		logger:  logrus.WithField("module", "vol_service").Logger,
		baseUrl: mockServer.URL,
	}
	expect := map[string]float64{
		"0x0b1a6fdd08b8e63d6b9476b971f03354823448ce": 173.7686,
	}
	res, err := li.processVolume(1730468849, 1735134449, "t")
	assert.NoErrorf(t, err, "Failed to get: %v", err)
	assert.Equal(t, expect["0x0b1a6fdd08b8e63d6b9476b971f03354823448ce"], res["0x0b1a6fdd08b8e63d6b9476b971f03354823448ce"])
}
