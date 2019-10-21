package exlru

import (
	"context"
	"fmt"
	"time"
)

type RunFunc func(ctx context.Context) (interface{}, error)

type MemCache struct {
	caches map[string]*ExCache
}

func NewMemCache() *MemCache {
	return &MemCache{
		caches: make(map[string]*ExCache),
	}
}


// Execute
// name lru cache instance name
// key cache key
// runFunc run func in real
// maxEntries max length
// expire expire time
func (m *MemCache) Execute(ctx context.Context, name string, key string, runFunc RunFunc, maxEntries int, expire *time.Duration) (interface{}, error) {
	cache, ok := m.caches[name]
	if !ok {
		cache = NewExLru(maxEntries)
		m.caches[name] = cache
	}
	dest, ok := cache.Get(key)
	if !ok {
		dest, err := runFunc(ctx)
		if err != nil {
			return nil, err
		}
		if expire != nil {
			cache.AddWithExpire(key, dest, *expire)
		} else {
			cache.Add(key, dest)
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
