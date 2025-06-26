package volume

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type midgardTracker struct {
	baseUrl        string
	chainDecimal   int
	clientIdHeader string
	logger         *logrus.Logger
}

func NewMidgardVolumeTracker(baseAddress string, chainDecimal int, clientIdHeader string) IVolumeTracker {
	return &midgardTracker{
		baseUrl:        fmt.Sprintf("%s/v2/actions", baseAddress),
		chainDecimal:   chainDecimal,
		clientIdHeader: clientIdHeader,
		logger:         logrus.WithField("module", "midgard_tracker").Logger,
	}
}

func (v *midgardTracker) SafeClose(closer io.Closer) {
	if err := closer.Close(); err != nil {
		v.logger.Error(err)
	}
}

func (v *midgardTracker) FetchVolume(from, to int64, affiliate string) (map[string]float64, error) {
	return v.processVolumeWithToken(from, to, affiliate, "")
}

func (v *midgardTracker) processVolumeWithToken(from, to int64, affiliate, nextPageToken string) (map[string]float64, error) {
	time.Sleep(1 * time.Second) // to avoid hitting rate limits
	url := fmt.Sprintf("%s?affiliate=%s&type=swap&timestamp=%d", v.baseUrl, affiliate, to)
	if nextPageToken != "" {
		url = fmt.Sprintf("%s?affiliate=%s&type=swap&nextPageToken=%s", v.baseUrl, affiliate, nextPageToken)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	if v.clientIdHeader != "" {
		req.Header.Set("X-Client-ID", v.clientIdHeader)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer v.SafeClose(resp.Body)

	if resp.StatusCode == http.StatusTooManyRequests {
		time.Sleep(30 * time.Second)
		return v.processVolumeWithToken(from, to, affiliate, nextPageToken)
	}
	var volRes tcVolumeModel
	if err := json.NewDecoder(resp.Body).Decode(&volRes); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	res := make(map[string]float64)
	for _, action := range volRes.Actions {
		if action.Status != "success" {
			continue
		}
		// convert nanoseconds to seconds
		date := action.Date / 1e9
		if date < from {
			return res, nil
		}
		for _, out := range action.Out {
			if out.Affiliate != nil && *out.Affiliate {
				continue
			}
			for _, outCoin := range out.OutCoins {
				res[out.Address] += float64(outCoin.Amount) * math.Pow10(-v.chainDecimal) * action.Metadata.Swap.OutPriceUSD
			}
		}
	}
	if volRes.Meta.NextPageToken != "" {
		nextRes, err := v.processVolumeWithToken(from, to, affiliate, volRes.Meta.NextPageToken)
		if err != nil {
			return nil, err
		}
		for k, v := range nextRes {
			res[k] += v
		}
	}
	return res, nil
}

type tcVolumeModel struct {
	Actions []tcActions `json:"actions"`
	Meta    tcMeta      `json:"meta"`
}
type tcSwap struct {
	OutPriceUSD float64 `json:"outPriceUSD,string"`
}
type tcMetadata struct {
	Swap tcSwap `json:"swap"`
}
type tcOutCoins struct {
	Amount int64  `json:"amount,string"`
	Asset  string `json:"asset"`
}
type tcOut struct {
	Address   string       `json:"address"`
	Affiliate *bool        `json:"affiliate"`
	OutCoins  []tcOutCoins `json:"coins"`
}
type tcActions struct {
	Date     int64      `json:"date,string"`
	Metadata tcMetadata `json:"metadata"`
	Out      []tcOut    `json:"out"`
	Status   string     `json:"status"`
}
type tcMeta struct {
	NextPageToken string `json:"nextPageToken"`
	PrevPageToken string `json:"prevPageToken"`
}
