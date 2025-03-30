package volume

import "io"

type IVolumeTracker interface {
	SafeClose(closer io.Closer)
	FetchVolume(from, to int64, affiliate string) (map[string]float64, error)
}
