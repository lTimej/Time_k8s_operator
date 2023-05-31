package caches

import (
	"Time_k8s_operator/internal/dao"
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/pkg/cache"
	"fmt"
	"strconv"
	"time"
)

type SpaceCache struct {
	cache *cache.Cache
}

func newSpaceCache() *SpaceCache {
	sc := &SpaceCache{
		cache: cache.NewCache("space-spe"),
	}
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			<-ticker.C
			sc.refresh()
		}
	}()
	return sc
}

func (sc *SpaceCache) refresh() {
	specs := dao.FindAllSpec()
	fmt.Println(specs, "$$$$$$$$$$$$$$$$$$$$$$$$")
	m := make(map[string]interface{}, len(specs))
	for i, _ := range specs {
		m[strconv.Itoa(int(specs[i].Id))] = &specs[i]
	}
	sc.cache.Replace(m)
}

func (sc *SpaceCache) LoadCache() {
	specs := dao.FindAllSpec()
	for i, _ := range specs {
		spec := specs[i]
		sc.cache.Set(strconv.Itoa(int(spec.Id)), &spec)
	}
}

func (sc *SpaceCache) GetSpaceSpec(key uint32) *model.SpaceSpec {
	val, ok := sc.cache.GetByInt(int(key))
	if !ok {
		return nil
	}
	item, ok := val.(*model.SpaceSpec)
	if !ok {
		return nil
	}
	p := *(item)
	return &p
}

func (sc *SpaceCache) GetAllSpaceSpec() (specs []*model.SpaceSpec) {
	alls := sc.cache.GetAll()
	for _, all := range alls {
		specs = append(specs, all.(*model.SpaceSpec))
	}
	return
}
