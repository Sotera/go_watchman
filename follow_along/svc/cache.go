package main

// cache is string-keyed cache with eviction by specified max cache size.
type cache struct {
	maxSize int
	items   []cacheItem
}

type cacheItem struct {
	key   string
	value interface{}
}

// add if not present. evict from tail on each successful add.
func (c *cache) add(citem cacheItem) {
	if _, hit := c.item(citem.key); !hit {
		if len(c.items) >= c.maxSize {
			c.items = c.items[:c.maxSize-1]
		}
		// prepends, by inverting append
		c.items = append([]cacheItem{citem}, c.items...)
	}
}

func (c *cache) item(key string) (citem cacheItem, hit bool) {
	for _, v := range c.items {
		if v.key == key {
			return v, true
		}
	}
	return cacheItem{}, false
}
