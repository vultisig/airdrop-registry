package volume

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
)

type lifiVolumeTracker struct {
	baseUrl string
	logger  *logrus.Logger
}

func NewLifiVolumeTracker() IVolumeTracker {
	return &lifiVolumeTracker{
		baseUrl: "https://li.quest/v1",
		logger:  logrus.WithField("module", "lifi_volume_tracker").Logger,
	}
}

func (l *lifiVolumeTracker) SafeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		l.logger.Error(err)
	}
}
func (l *lifiVolumeTracker) FetchVolume(from, to int64, affiliate string) (map[string]float64, error) {
	res := make(map[string]float64)
	if !l.isValidAffiliate(affiliate) {
		return res, nil
	}
	url := fmt.Sprintf("%s/analytics/transfers?integrator=%s&fromTimestamp=%d&toTimestamp=%d", l.baseUrl, affiliate, from, to)
	resp, err := http.Get(url)
	if err != nil {
		l.logger.WithError(err).Error("error making GET request")
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer l.SafeClose(resp.Body)
	if resp.StatusCode != http.StatusOK {
		l.logger.WithField("status_code", resp.StatusCode).Error("error response from server")
		return nil, fmt.Errorf("error response from server: %s", resp.Status)
	}
	var volRes lifiVolumeModel
	if err := json.NewDecoder(resp.Body).Decode(&volRes); err != nil {
		l.logger.WithError(err).Error("error decoding response")
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	for _, transfer := range volRes.Transfers {
		if transfer.Status != "DONE" {
			l.logger.WithField("status", transfer.Status).Info("transfer not done")
			continue
		}
		res[transfer.ToAddress] += transfer.Receiving.AmountUSD
	}
	return res, nil
}

// affiliate must be an alphanumeric string (additionally - ,_, and . are allowed) with a maximum length of 23 characters
var affiliateRegex = regexp.MustCompile(`^[a-zA-Z0-9-_.]{1,23}$`)

func (l *lifiVolumeTracker) isValidAffiliate(affiliate string) bool {
	return affiliateRegex.MatchString(affiliate)
}

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
