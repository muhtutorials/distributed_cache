package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(key)] = value

	log.Printf("SET %s to %s\n", key, value)

	if ttl > 0 {
		go func() {
			time.Sleep(ttl)
			delete(c.data, string(key))
			log.Printf("(%s) key has been deleted\n", key)
		}()
	}

	return nil
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keyStr := string(key)

	val, ok := c.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", keyStr)
	}

	log.Printf("GET %s=%s\n", key, val)

	return val, nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(key)]

	return ok
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.data, string(key))

	return nil
}
