package utils

import "crypto/sha256"

// return sha256 hash of input
func SHA256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}
