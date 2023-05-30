package caches

import (
	"reflect"
	"sync"
)

type CacheFactory struct {
	caches map[reflect.Type]interface{}
	lock   sync.Mutex
}

func NewCacheFactory() *CacheFactory {
	return &CacheFactory{
		caches: make(map[reflect.Type]interface{}),
	}
}

func (cf *CacheFactory) TemplateCache() *TemplateCache {
	t := reflect.TypeOf(&TemplateCache{})
	cf.lock.Lock()
	defer cf.lock.Unlock()
	if c, ok := cf.caches[t]; ok {
		return c.(*TemplateCache)
	}
	temp := newTemplateCache()
	cf.caches[t] = temp
	temp.LoadCache()
	return temp
}

func (cf *CacheFactory) SpaceSpecCache() *SpaceCache {
	t := reflect.TypeOf(&SpaceCache{})
	cf.lock.Lock()
	defer cf.lock.Unlock()
	if c, ok := cf.caches[t]; ok {
		return c.(*SpaceCache)
	}
	space := newSpaceCache()
	cf.caches[t] = space
	space.LoadCache()
	return space
}
