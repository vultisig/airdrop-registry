package volume

import "io"

type IVolume interface {
	closer(closer io.Closer)
	processVolume(from, to int64, affiliate string) (map[string]float64, error)
}
