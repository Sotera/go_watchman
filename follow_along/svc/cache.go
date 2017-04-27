package main

import "sync"

// cache is string-keyed cache with eviction by specified max cache size.
// thread-safe.
type cache struct {
	maxSize int
	items   []cacheItem
	sync.RWMutex
}

type cacheItem struct {
	key   string
	value interface{}
}

// add if not present. evict from tail on each successful add.
func (c *cache) add(citem cacheItem) {
	if c.maxSize == 0 {
		panic("max size not set")
	}
	c.Lock()
	defer c.Unlock()

	if _, hit := c._item(citem.key); !hit {
		if len(c.items) >= c.maxSize {
			c.items = c.items[:c.maxSize-1]
		}
		// prepends, by inverting append
		c.items = append([]cacheItem{citem}, c.items...)
	}
}

// safely get item.
func (c *cache) get(key string) (citem cacheItem, hit bool) {
	c.Lock()
	defer c.Unlock()
	return c._item(key)
}

// unsafely get item. for internal use only.
func (c *cache) _item(key string) (citem cacheItem, hit bool) {
	for _, v := range c.items {
		if v.key == key {
			return v, true
		}
	}
	return cacheItem{}, false
}

func (c *cache) clear() {
	c.Lock()
	defer c.Unlock()
	c.items = []cacheItem{}
}
