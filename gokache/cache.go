package gokache

import (
	"gokache/lru"
	"sync"
)

type cache struct {
	mu       sync.Mutex
	lru      *lru.Lru
	capacity int64
}

func (c *cache) set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = lru.InitLru(c.capacity, 0, nil)
	}

	c.lru.Set(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// concurrence test
	//fmt.Println("----------")
	//fmt.Println(key)
	//fmt.Println(key)
	//time.Sleep(1 * time.Second)
	//fmt.Println(key)
	//fmt.Println(key)
	//fmt.Println("----------")

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		// return "val,true" if v belong to ByteView, otherwise, return ",false"
		return v.(ByteView), ok
	}

	return
}
