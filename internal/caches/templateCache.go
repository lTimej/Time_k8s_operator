package caches

import (
	"Time_k8s_operator/internal/dao"
	"Time_k8s_operator/internal/model"
	"Time_k8s_operator/pkg/cache"
	"fmt"
	"time"
)

const (
	TemplateKey = "templates"
	KindsKey    = "kinds"
)

type TemplateCache struct {
	cache *cache.Cache
}

func newTemplateCache() *TemplateCache {
	t := &TemplateCache{
		cache: cache.NewCache("space-template"),
	}
	go func() {
		tick := time.NewTicker(time.Minute)
		defer tick.Stop()
		for {
			<-tick.C
			t.LoadCache()
		}
	}()
	return t
}

func (tc *TemplateCache) LoadCache() {
	tmps := dao.FindAllSpaceTemplateByUsing()
	tps := make(map[uint32]*model.SpaceTemplate, len(tmps))
	for i := 0; i < len(tmps); i++ {
		tmp := tmps[i]
		tps[tmp.Id] = &tmp
	}
	fmt.Printf("tps:%#v\n", tps)
	tc.cache.Set(TemplateKey, tps)

	kinds := dao.FindAllTemplateKind()
	kis := make(map[uint32]*model.TemplateKind, len(kinds))
	for _, kind := range kinds {
		ki := kind.Id
		kis[ki] = &kind
	}
	tc.cache.Set(KindsKey, kis)
}

func (tc *TemplateCache) GetSpaceTemplate(st_id uint32) *model.SpaceTemplate {
	item, ok := tc.cache.Get(TemplateKey)
	if !ok {
		return nil
	}
	templates, ok := item.(map[uint32]*model.SpaceTemplate)
	if !ok {
		return nil
	}
	temp := model.SpaceTemplate{}
	if template, ok := templates[st_id]; ok {
		temp = *template
	} else {
		return nil
	}

	return &temp
}

func (tc *TemplateCache) GetAllSpaceTemplate() (templates []*model.SpaceTemplate) {
	item, ok := tc.cache.Get(TemplateKey)
	if !ok {
		return nil
	}
	template, ok := item.(map[uint32]*model.SpaceTemplate)
	if !ok {
		return nil
	}
	i := 0
	for _, val := range template {
		p := *(val)
		templates[i] = &p
		i += 1
	}
	return
}

func (tc *TemplateCache) GetAllKind() (kinds []*model.TemplateKind) {
	item, ok := tc.cache.Get(KindsKey)
	if !ok {
		return nil
	}
	kind, ok := item.(map[uint32]*model.TemplateKind)
	if !ok {
		return nil
	}
	i := 0
	for _, val := range kind {
		k := *(val)
		kinds[i] = &k
		i += 1
	}
	return
}
