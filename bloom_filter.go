package bloom_filter

import (
	"sync"
)

const defaultBucketLen = 1 << 16

var defaultHashFuncs = []BaseHashFunc{
	md5HashFunc,
	sha256HashFunc,
}

type bucket []bool

func newBucket(len int64) bucket {
	return make([]bool, len)
}

func (b *bucket) add(n int64) {
	modN := n % int64(len(*b))
	i := len(*b) - 1
	for modN > 0 {
		[]bool(*b)[i] = (modN % 2) == 1
		modN /= 2
		i--
	}
}

func (b *bucket) contains(n int64) bool {
	modN := n % int64(len(*b))
	i := len(*b) - 1
	for modN > 0 {
		if []bool(*b)[i] != ((modN % 2) == 1) {
			return false
		}
		modN /= 2
		i--
	}

	return true
}

type BloomFilter struct {
	mu        sync.RWMutex
	bucket    bucket
	hashFuncs []BaseHashFunc
}

func DefaultBloomFilter() *BloomFilter {
	bf := new(BloomFilter)
	bf.hashFuncs = defaultHashFuncs
	bf.bucket = newBucket(defaultBucketLen)
	return bf
}

func NewBloomFilter(bucketLen int64, hashFuncs []BaseHashFunc) *BloomFilter {
	bf := new(BloomFilter)
	bf.hashFuncs = hashFuncs
	bf.bucket = newBucket(bucketLen)
	return bf
}

func (bf *BloomFilter) Add(data []byte) {
	bf.mu.Lock()
	defer bf.mu.Lock()
	for _, f := range bf.hashFuncs {
		bf.bucket.add(f(data))
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	bf.mu.RLock()
	defer bf.mu.Unlock()
	for _, f := range bf.hashFuncs {
		if !bf.bucket.contains(f(data)) {
			return false
		}
	}

	return true
}

func (bf *BloomFilter) Cap() int64 {
	bf.mu.RLock()
	defer bf.mu.Unlock()
	return int64(len(bf.bucket))
}
