package cache

import (
	"fmt"
	"strconv"
	"sync"
)

type Cache struct {
	lock sync.RWMutex
	data map[string]interface{}
}

var (
	caches = map[string]*Cache{}
	lock   sync.Mutex
)

func NewCache(name string) *Cache {
	lock.Lock()
	defer lock.Unlock()
	if c, ok := caches[name]; ok {
		return c
	}
	c := &Cache{
		data: make(map[string]interface{}),
	}
	caches[name] = c
	return c
}

func (c *Cache) Set(name string, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[name] = val
}

func (c *Cache) Get(name string) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	val, ok := c.data[name]
	return val, ok
}

func (c *Cache) GetByInt(key int) (interface{}, bool) {
	return c.Get(strconv.Itoa(key))
}

func (c *Cache) GetAll() []interface{} {
	c.lock.RLock()
	fmt.Println(111111111)
	defer c.lock.RUnlock()
	fmt.Println(222222)
	ret := make([]interface{}, 0, len(c.data))
	for _, val := range c.data {
		ret = append(ret, val)
	}
	return ret
}

func (c *Cache) Replace(data map[string]interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data = data
}
