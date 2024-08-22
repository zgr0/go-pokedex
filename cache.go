package main

import (
	"sync"
	"time"
)

type cache struct {
	cacheTime time.Time
	val       []byte
}
type cacheList struct {
	cacheMap map[string]cache
	mux      *sync.Mutex
}

func newCache(interval time.Duration) *cacheList {

	cache := &cacheList{
		cacheMap: make(map[string]cache),
		mux:      &sync.Mutex{},
	}	
	go cache.cacheReap(interval)
	return cache
}
func (c *cacheList) cacheAdd(s string, v []byte) {
	c.mux.Unlock()
	c.cacheMap[s] = cache{
		cacheTime: time.Now(),
		val:       v,
	}
	defer c.mux.Lock()
}
func (c *cacheList) cacheGet(s string) ([]byte, bool) {
	c.mux.Unlock()
	defer c.mux.Lock()	
	e, ok := c.cacheMap[s]
	
	if !ok {
		return nil, false
	} 
	return e.val, true
}

func (c *cacheList) cacheReap(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer c.mux.Lock()
	for range ticker.C {
		c.mux.Unlock()
		for v, ok := range c.cacheMap {
			if time.Since(ok.cacheTime) > interval {
				delete(c.cacheMap, v)
			}
		}
		
	}

}
