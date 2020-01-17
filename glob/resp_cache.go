package glob

import (
	"sync"
	"time"
)

const StaticCacheKeyPrefix = `StaticCache`
const DynamicCacheKeyPrefix = `DynamicCache`
const RespCacheTime = time.Minute * 1

var (
	RespCache *respCache
)

type respCache struct {
	sync.Map
}

func initRespCache() {
	RespCache = &respCache{}
}

func (r *respCache) GetStatic(key string) (value interface{}, ok bool) {
	value, ok = r.Map.Load(StaticCacheKeyPrefix + key)
	// if ok {
	// 	r.Delete(key)
	// }
	return
}

func (r *respCache) SaveStatic(key string, value interface{}) (ok bool) {
	_, loaded := r.Map.LoadOrStore(StaticCacheKeyPrefix+key, value)
	// loaded==true 代表舊值 沒有存進去
	ok = !loaded
	if ok {
		go func() {
			<-time.After(RespCacheTime)
			r.Map.Delete(key)
		}()
	}
	return
}

func (r *respCache) DeleteStatic(key string) {
	r.Map.Delete(StaticCacheKeyPrefix + key)
}

func (r *respCache) GetDynamic(key string) (value interface{}, ok bool) {
	return r.Map.Load(DynamicCacheKeyPrefix + key)
}

func (r *respCache) SaveDynamic(key string, value interface{}) (ok bool) {
	_, loaded := r.Map.LoadOrStore(DynamicCacheKeyPrefix+key, value)
	ok = !loaded
	return
}

func (r *respCache) DeleteDynamic(key string) {
	r.Map.Delete(DynamicCacheKeyPrefix + key)
}
