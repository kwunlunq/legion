package service

import (
	"errors"
	"sync"

	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

var (
	caches = cache{
		staticCaches:  make(map[string][]byte),
		dynamicCaches: make(map[string][]byte),
	}
)

type cache struct {
	sync.Mutex
	staticCaches  map[string][]byte
	dynamicCaches map[string][]byte
}

func GetStaticCaches(req model.CacheRequest) (*model.CacheResponse, error) {
	resp := &model.CacheResponse{}

	val, ok := caches.staticCaches[req.TaskID]
	if !ok {
		return nil, errors.New("task_id:" + req.TaskID + " 不存在")
	}

	resp.TaskID = req.TaskID
	resp.Content = string(val)

	caches.Lock()
	delete(caches.staticCaches, resp.TaskID)
	caches.Unlock()

	return resp, nil
}

func GetDynamicCaches(req model.CacheRequest) (*model.CacheResponse, error) {
	resp := &model.CacheResponse{}

	val, ok := caches.dynamicCaches[req.TaskID]
	if !ok {
		return nil, errors.New("task_id:" + req.TaskID + " 不存在")
	}

	resp.TaskID = req.TaskID
	resp.Content = string(val)

	caches.Lock()
	delete(caches.dynamicCaches, resp.TaskID)
	caches.Unlock()

	return resp, nil
}
