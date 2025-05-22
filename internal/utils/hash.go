package utils

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

// return sha256 hash of input
func SHA256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}

func Hash160(data []byte) []byte {
	ripemd := ripemd160.New()
	ripemd.Write(SHA256(data)[:])
	return ripemd.Sum(nil)
}
