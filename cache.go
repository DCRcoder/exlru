package exlru

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RunFunc func(ctx context.Context) (interface{}, error)

type MemCache struct {
	caches *sync.Map
	mux    *sync.Mutex
}

func NewMemCache() *MemCache {
	return &MemCache{
		caches: &sync.Map{},
		mux:    &sync.Mutex{},
	}
}

// Execute try get result from memcache else run call function and set cache
// name lru cache instance name
// key cache key
// runFunc run func in real
// maxEntries max length
// expire expire time
func (m *MemCache) Execute(ctx context.Context, name string, key string, runFunc RunFunc, maxEntries int, expire *time.Duration) (interface{}, error) {
	m.mux.Lock()
	cache, ok := m.caches.Load(name)
	if !ok {
		cache = NewExLru(maxEntries)
		m.caches.Store(name, cache)
	}
	m.mux.Unlock()
	c := cache.(*ExCache)
	dest, ok := c.Get(key)
	if !ok {
		dest, err := runFunc(ctx)
		if err != nil {
			return nil, err
		}
		if expire != nil {
			c.AddWithExpire(key, dest, *expire)
		} else {
			c.Add(key, dest)
		}
		return dest, nil
	}
	return dest, nil
}

func GenKey(param ...interface{}) string {
	key := ""
	for _, p := range param {
		key += fmt.Sprintf("%v", p) + "|"
	}
	return key
}
