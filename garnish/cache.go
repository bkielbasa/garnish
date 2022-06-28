package garnish

import (
	"sync"
	"time"
)

type data struct {
	data    []byte
	expires *time.Time
}

func shouldBeCleared(d data) bool {
	if d.expires == nil {
		return false
	}

	return time.Now().After(*d.expires)
}

type cache struct {
	mutex *sync.Mutex
	data  map[string]data
}

//The cache keeps the data in a map
func (c *cache) store(key string, rawData []byte, timeout time.Duration) {
	d := data{
		data: rawData,
	}
	if timeout != 0 {
		t := time.Now().Add(timeout)
		d.expires = &t
	}

	c.mutex.Lock()
	c.data[key] = d
	c.mutex.Unlock()

	time.AfterFunc(timeout, func() {
		c.clear(key)
	})
}

func (c *cache) clear(key string) {
	c.mutex.Lock()
	c.clearKey(key)
	c.mutex.Unlock()
}

func (c *cache) clearKey(key string) {
	delete(c.data, key)
}

func (c cache) get(key string) []byte {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if d, ok := c.data[key]; ok {
		if shouldBeCleared(d) {
			c.clearKey(key)
			return nil
		}

		return d.data
	}
	return nil
}

func newCache() *cache {
	return &cache{data: map[string]data{}, mutex: &sync.Mutex{}}
}
