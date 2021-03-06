// +build ignore

package main

import (
	"strconv"
	"sync"
)

// check for races with go run:
// go run -race race.go
func main() {
	c := cache{maxSize: 100}
	// writes
	for i := 0; i < 100; i++ {
		go func(n int) {
			c.add(cacheItem{key: strconv.Itoa(n), value: ""})
		}(i)
	}
	// reads
	for i := 0; i < 100; i++ {
		go func(n int) {
			c.get(strconv.Itoa(n))
		}(i)
	}
}

/////////////////////////////////
// copied cache.go content below.
/////////////////////////////////
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
