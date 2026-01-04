package cache

import "github.com/dgraph-io/ristretto"

func NewRistrettoCache() (*ristretto.Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     200 << 20, // 200 MB
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}

	return cache, nil
}
