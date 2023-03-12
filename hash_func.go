package bloom_filter

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
)

type BaseHashFunc func([]byte) int64

func md5HashFunc(data []byte) int64 {
	h := md5.New()
	return bytesToInt64(h.Sum(data))
}

func sha256HashFunc(data []byte) int64 {
	h := sha256.New()
	return bytesToInt64(h.Sum(data))
}

func bytesToInt64(data []byte) int64 {
	return int64(binary.BigEndian.Uint64(data))
}
