package service

import (
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

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

	if _, ok := caches.staticCaches[req.TaskID]; !ok {
		caches.Lock()
		caches.staticCaches[req.TaskID] = body
		caches.Unlock()
	}

	return resp, nil
}
