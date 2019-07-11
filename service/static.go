package service

import (
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

func StaticScrape(req *model.Request) error {
	doc, err := glob.GetAndConvertToDocument(req.URL)

	body := []byte(doc.Find(req.Target).Text())

	if req.Charset != "" {
		body, err = glob.Decoder(body, req.Charset)
	}

	req.Body = body
	req.Error = err

	if err := glob.Cache.SetStaticCache(req.TaskID, body); err != nil {
		return err
	}

	return nil
}

func GetStaticCache(req *model.CacheRequest) error {
	value, err := glob.Cache.GetStaticCache(req.TaskID)
	if err != nil {
		return err
	}

	req.Content = string(value)

	return nil
}
