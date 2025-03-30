package volume

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

//go:embed tc_vol_one.json
var tcVolumeResponseOne string

func TestTCVolume(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response2 := map[string]any{
			"actions": []any{},
			"meta": map[string]any{
				"nextPageToken": "",
				"prevPageToken": "",
			},
		}
		if r.URL.Query().Get("nextPageToken") == "198620799000000001" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response2)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(tcVolumeResponseOne))
		}
	}))
	defer mockServer.Close()
	vr := tcVolumeTrack{
		logger:  logrus.WithField("module", "tc_vol_service").Logger,
		baseUrl: mockServer.URL,
	}
	expect := map[string]float64{
		"0x060c27cd6719477f233e403d74da9513886f0a1a": 324578205539.9229,
	}
	res, err := vr.processVolume(1739510000, 1739519656, "t")
	assert.NoErrorf(t, err, "Failed to get: %v", err)
	assert.Equal(t, expect, res)
}
