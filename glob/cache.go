package glob

import (
	"errors"
	"sync"
)

type cache struct {
	sync.Mutex
	staticCache  map[string][]byte
	dynamicCache map[string][]byte
}

func NewCache() *cache {
	c := &cache{
		staticCache:  make(map[string][]byte),
		dynamicCache: make(map[string][]byte),
	}

	return c
}

func (c *cache) GetStaticCache(key string) ([]byte, error) {
	val, ok := c.staticCache[key]
	if !ok {
		return nil, errors.New("key: " + key + " 不存在")
	}

	c.Lock()
	delete(c.staticCache, key)
	c.Unlock()

	return val, nil
}

func (c *cache) GetDynamicCache(key string) ([]byte, error) {
	val, ok := c.dynamicCache[key]
	if !ok {
		return nil, errors.New("key: " + key + " 不存在")
	}

	c.Lock()
	delete(c.dynamicCache, key)
	c.Unlock()

	return val, nil
}

func (c *cache) SetStaticCache(key string, value []byte) error {
	_, ok := c.staticCache[key]
	if ok {
		return errors.New("key: " + key + " 已存在")
	}

	c.Lock()
	c.staticCache[key] = value
	c.Unlock()

	return nil
}

func (c *cache) SetDynamicCache(key string, value []byte) error {
	_, ok := c.dynamicCache[key]
	if ok {
		return errors.New("key: " + key + " 已存在")
	}

	c.Lock()
	c.dynamicCache[key] = value
	c.Unlock()

	return nil
}
