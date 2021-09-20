package lru

import "github.com/hashicorp/golang-lru"

type LRU struct {
	cache *lru.Cache
}

func New(size int) (*LRU, error) {
	hashicorpLRU, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	return &LRU{cache: hashicorpLRU}, nil
}

func (c *LRU) Set(key string, data []byte) {
	c.cache.Add(key, data)
}

func (c *LRU) Get(key string) ([]byte, bool) {
	data, ok := c.cache.Get(key)
	if !ok {
		return nil, ok
	}

	bytes, ok := data.([]byte)
	return bytes, ok
}

func (c *LRU) Del(key string) {
	c.cache.Remove(key)
}
