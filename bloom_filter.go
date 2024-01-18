package bloom_filter

import (
	"sync"
)

const defaultBucketLen = 1 << 16

var defaultHashFuncs = []BaseHashFunc{
	md5HashFunc,
	sha256HashFunc,
}

type BloomFilter struct {
	mu        sync.RWMutex
	buckets   buckets
	hashFuncs []BaseHashFunc
}

func DefaultBloomFilter() *BloomFilter {
	bf := new(BloomFilter)
	bf.hashFuncs = defaultHashFuncs
	bf.buckets = newBucket(defaultBucketLen)
	return bf
}

func NewBloomFilter(bucketLen int64, hashFuncs []BaseHashFunc) *BloomFilter {
	bf := new(BloomFilter)
	bf.hashFuncs = hashFuncs
	bf.buckets = newBucket(bucketLen)
	return bf
}

func (bf *BloomFilter) Add(data []byte) {
	bf.mu.Lock()
	defer bf.mu.Lock()
	for _, f := range bf.hashFuncs {
		bf.buckets.add(f(data))
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	bf.mu.RLock()
	defer bf.mu.Unlock()
	for _, f := range bf.hashFuncs {
		if !bf.buckets.contains(f(data)) {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) Cap() int64 {
	bf.mu.RLock()
	defer bf.mu.Unlock()
	return int64(len(bf.buckets))
}
