package cache

import (
	"sync"
)

type InMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat
}

func (c *InMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp, exist := c.c[k]
	if exist {
		c.del(k, tmp)
	}
	c.c[k] = v
	c.add(k, v)
	return nil
}

func (c *InMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k], nil
}

func (c *InMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if exist {
		delete(c.c, k)
		c.del(k, v)
	}
	return nil
}

func (c *InMemoryCache) GetStat() Stat {
	return c.Stat
}

func newInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		c:     make(map[string][]byte),
		mutex: sync.RWMutex{},
		Stat:  Stat{},
	}
}

func (c *InMemoryCache) NewScanner() Scanner {
	pairCh := make(chan *pair)
	closeCh := make(chan *struct{})
	go func() {
		defer close(closeCh)
		c.mutex.RLock()
		for k, v := range c.c {
			c.mutex.RUnlock()
			select {
			case <-closeCh:
				return
			case pairCh <- &pair{k, v}:
			}
			c.mutex.RLock()
		}
	}()
}
