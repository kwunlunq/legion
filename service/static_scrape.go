package service

import (
	"errors"
	"sync"

	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

var (
	staticCaches = cache{
		caches: make(map[string][]byte),
	}
)

type cache struct {
	sync.Mutex
	caches map[string][]byte
}

func StaticScrape(req model.Request) (*model.Response, error) {
	doc, err := glob.GetAndConvertToDocument(req.URL)

	body := []byte(doc.Find(req.Target).Text())

	if req.Charset != "" {
		body, err = glob.Decoder(body, req.Charset)
	}

	resp := &model.Response{}
	resp.TaskID = req.TaskID
	resp.Body = body
	resp.Error = err

	if _, ok := staticCaches.caches[req.TaskID]; !ok {
		staticCaches.Lock()
		staticCaches.caches[req.TaskID] = body
		staticCaches.Unlock()
	}

	return resp, nil
}

func GetStaticCache(req model.StaticRequest) (*model.StaticResponse, error) {
	resp := &model.StaticResponse{}

	val, ok := staticCaches.caches[req.TaskID]
	if !ok {
		return nil, errors.New("task_id:" + req.TaskID + " 不存在")
	}

	resp.TaskID = req.TaskID
	resp.Content = string(val)

	staticCaches.Lock()
	delete(staticCaches.caches, resp.TaskID)
	staticCaches.Unlock()

	return resp, nil
}
