package cache

import (
	"errors"
	"sync"

	"github.com/Visharad18/feedadapter/entity"
)

var (
	errInvalidValueType = errors.New("value of invalid type supplied")
	errKeyNotFound      = errors.New("key not found in cache")
)

type InMemoryCache struct {
	cache map[string]map[int64]*entity.HistoricalData // map[symbol][unixEpoch]{OHLCV}
	mu    sync.RWMutex
}

func NewInMeInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]map[int64]*entity.HistoricalData),
		mu:    sync.RWMutex{},
	}
}

func (c *InMemoryCache) Store(key string, value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := value.(map[int64]*entity.HistoricalData)
	if !ok {
		return errInvalidValueType
	}

	if _, ok := c.cache[key]; !ok {
		c.cache[key] = make(map[int64]*entity.HistoricalData)
	}

	for t, hd := range val {
		c.cache[key][t] = hd
	}

	return nil
}

func (c *InMemoryCache) Get(key string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.cache[key]
	if !ok {
		return nil, errKeyNotFound
	}

	return value, nil
}

func (c *InMemoryCache) GetAll() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	res := make(map[string]any)

	for k, v := range c.cache {
		res[k] = v
	}

	return res
}
