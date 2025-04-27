package volume

import (
	"io"

	"github.com/vultisig/airdrop-registry/config"
)

type IVolumeTracker interface {
	SafeClose(closer io.Closer)
	FetchVolume(from, to int64, affiliate string) (map[string]float64, error)
}

type VolumeResolver struct {
	trackers  []IVolumeTracker
	volume    map[string]float64
	affiliate []string
}

func NewVolumeResolver(cfg *config.Config) (*VolumeResolver, error) {
	pr := &VolumeResolver{
		affiliate: cfg.VolumeTrackingAPI.AffiliateAddress,
		trackers: []IVolumeTracker{
			NewMidgardVolumeTracker(cfg.VolumeTrackingAPI.TCMidgardBaseURL, 8),
			NewMidgardVolumeTracker(cfg.VolumeTrackingAPI.MayaMidgardBaseURL, 10),
			NewLifiVolumeTracker(),
			NewOneInchVolumeTracker(cfg.VolumeTrackingAPI.EtherscanAPIKey, cfg.VolumeTrackingAPI.EthplorerAPIKey),
		},
		volume: make(map[string]float64),
	}
	return pr, nil
}

func (v *VolumeResolver) LoadVolume(from, to int64) error {
	res := make(map[string]float64)
	for _, aff := range v.affiliate {
		for _, tracker := range v.trackers {
			vol, err := tracker.FetchVolume(from, to, aff)
			if err != nil {
				return err
			}
			for k, v := range vol {
				res[k] += v
			}
		}
	}
	v.volume = res
	return nil
}

func (v *VolumeResolver) GetVolume(address string) float64 {
	return v.volume[address]
}
