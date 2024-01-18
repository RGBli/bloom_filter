package bloom_filter

type buckets []bool

func newBucket(len int64) buckets {
	return make([]bool, len)
}

func (b *buckets) add(n int64) {
	modN := n % int64(len(*b))
	i := len(*b) - 1
	for modN > 0 {
		[]bool(*b)[i] = (modN % 2) == 1
		modN /= 2
		i--
	}
}

func (b *buckets) contains(n int64) bool {
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
