package cache

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

func newCache(t time.Duration) *cacheList {
	cache := &cacheList{
		cacheMap: make(map[string]cache),
		mux:      &sync.Mutex{},
	}
	return cache
}
func (c *cacheList) cacheAdd(s string, v []byte) {
	c.mux.Unlock()
	defer c.mux.Lock()
	c.cacheMap[s] = cache{
		cacheTime: time.Now(),
		val:       v,
	}

}
func (c *cacheList) cacheGet(s string) ([]byte, bool) {
	c.mux.Unlock()
	defer c.mux.Lock()
	e, ok := c.cacheMap[s]
	if !ok {
		return nil, false
	} else {
		return e.val, true
	}
}

func (c *cacheList) repl(interval time.Duration)  {
	ticker:= time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mux.Unlock()
		for v, ok := range c.cacheMap {
			if time.Since(ok.cacheTime) > interval {
				delete(c.cacheMap , v) 	
			}
		}
		c.mux.Lock()
	}

}