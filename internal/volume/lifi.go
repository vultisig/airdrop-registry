package volume

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type lifiTransaction struct {
	AmountUSD float64 `json:"amountUSD,string"`
}
type lifiTransfer struct {
	Receiving lifiTransaction `json:"receiving"`
	ToAddress string          `json:"toAddress"`
	Status    string          `json:"status"`
}
type lifiVolumeModel struct {
	Transfers []lifiTransfer `json:"transfers"`
}
type lifiVolumeTrack struct {
	baseUrl string
	logger  *logrus.Logger
}

func NewVolumeTrack() IVolume {
	return &lifiVolumeTrack{
		baseUrl: "https://li.quest/v1",
		logger:  logrus.WithField("module", "lifi_volume_tracker").Logger,
	}
}

func (l *lifiVolumeTrack) closer(closer io.Closer) {
	if err := closer.Close(); err != nil {
		l.logger.Error(err)
	}
}
func (l *lifiVolumeTrack) processVolume(from, to int64, affiliate string) (map[string]float64, error) {
	url := fmt.Sprintf("%s/analytics/transfers?integrator=%s&fromTimestamp=%d&toTimestamp=%d", l.baseUrl, affiliate, from, to)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer l.closer(resp.Body)
	var volRes lifiVolumeModel
	if err := json.NewDecoder(resp.Body).Decode(&volRes); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	res := make(map[string]float64)
	for _, transfer := range volRes.Transfers {
		if transfer.Status != "DONE" {
			continue
		}
		res[transfer.ToAddress] += transfer.Receiving.AmountUSD
	}
	return res, nil
}
