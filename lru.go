package exlru

import (
	"container/list"
	"sync"
	"time"
)

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

type ExCache struct {
	MaxEntries int

	OnEvicted func(key Key, value interface{})

	ll *list.List

	cache map[interface{}]*list.Element

	mut sync.Mutex
}

type entry struct {
	key    Key
	value  interface{}
	expire int64
}

func NewExLru(maxEntries int) *ExCache {
	return &ExCache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

// Add adds a value to the cache.
func (c *ExCache) Add(key Key, value interface{}) {
	c.mut.Lock()
	defer c.mut.Unlock()
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}

	if ee, ok := c.cache[key]; ok && ee.Value.(*entry).expire > int64(time.Now().UnixNano()) {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
	}

	ele := c.ll.PushFront(&entry{key, value, 0})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// Add adds a value to the cache and set expire time.
func (c *ExCache) AddWithExpire(key Key, value interface{}, expire time.Duration) {
	c.mut.Lock()
	defer c.mut.Unlock()
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}

	if ee, ok := c.cache[key]; ok && ee.Value.(*entry).expire > int64(time.Now().UnixNano()) {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		ee.Value.(*entry).expire = int64(time.Now().UnixNano()) + expire.Nanoseconds()
	}

	ele := c.ll.PushFront(&entry{key, value, int64(time.Now().UnixNano()) + expire.Nanoseconds()})
	c.cache[key] = ele
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *ExCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	ele, ok := c.cache[key]
	if ok {
		deadline := ele.Value.(*entry).expire
		if deadline == 0 {
			c.ll.MoveToFront(ele)
			return ele.Value.(*entry).value, true
		}
		if deadline >= int64(time.Now().UnixNano()) {
			c.ll.MoveToFront(ele)
			return ele.Value.(*entry).value, true
		} else if deadline < int64(time.Now().UnixNano()) {
			c.ll.Remove(ele)
			return nil, false
		}
	}

	return
}

// Len returns the number of items in the cache.
func (c *ExCache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Remove removes the provided key from the cache.
func (c *ExCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *ExCache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *ExCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value)
	}
}

// Clear purges all stored items from the cache.
func (c *ExCache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
}
