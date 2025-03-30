package volume

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type midgardTracker struct {
	baseUrl string
	logger  *logrus.Logger
}

func NewTCMidgardTracker() IVolumeTracker {
	return &midgardTracker{
		baseUrl: "https://midgard.ninerealms.com/v2/actions",
		logger:  logrus.WithField("module", "thorchain_volume_tracker").Logger,
	}
}

func NewMayaMidgardTracker() IVolumeTracker {
	return &midgardTracker{
		baseUrl: "https://midgard.mayachain.info/v2/actions",
		logger:  logrus.WithField("module", "mayachain_volume_tracker").Logger,
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
	url := fmt.Sprintf("%s?affiliate=%s&type=swap&timestamp=%d", v.baseUrl, affiliate, to)
	if nextPageToken != "" {
		url = fmt.Sprintf("%s?affiliate=%s&type=swap&nextPageToken=%s", v.baseUrl, affiliate, nextPageToken)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer v.SafeClose(resp.Body)
	var volRes tcVolumeModel
	if err := json.NewDecoder(resp.Body).Decode(&volRes); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	res := make(map[string]float64)
	for _, action := range volRes.Actions {
		if action.Status != "success" {
			continue
		}
		date, err := strconv.ParseInt(action.Date[:10], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing action date: %w", err)
		}
		if date < from {
			return res, nil
		}
		for _, out := range action.Out {
			if out.Affiliate != nil && *out.Affiliate {
				continue
			}
			for _, outCoin := range out.OutCoins {
				amount, err := strconv.ParseInt(outCoin.Amount, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("error parsing out coin amount: %w", err)
				}
				res[out.Address] += float64(amount) * action.Metadata.Swap.OutPriceUSD
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
	Amount string `json:"amount"`
	Asset  string `json:"asset"`
}
type tcOut struct {
	Address   string       `json:"address"`
	Affiliate *bool        `json:"affiliate"`
	OutCoins  []tcOutCoins `json:"coins"`
}
type tcActions struct {
	Date     string     `json:"date"`
	Metadata tcMetadata `json:"metadata"`
	Out      []tcOut    `json:"out"`
	Status   string     `json:"status"`
}
type tcMeta struct {
	NextPageToken string `json:"nextPageToken"`
	PrevPageToken string `json:"prevPageToken"`
}
