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

	if err := glob.Cache.SetStaticCache(req.TaskID, body); err != nil {
		return nil, err
	}

	return resp, nil
}

func GetStaticCache(req model.CacheRequest) (*model.CacheResponse, error) {
	resp := &model.CacheResponse{}

	resp.TaskID = req.TaskID

	value, err := glob.Cache.GetStaticCache(req.TaskID)
	if err != nil {
		return nil, err
	}

	resp.Content = string(value)

	return resp, nil
}
