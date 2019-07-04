package service

import (
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

func StaticScrape(req model.Request) (*model.Response, error) {
	doc, err := glob.GetAndConvertToDocument(req.URL)
	if err != nil {
		return nil, err
	}

	body := []byte(doc.Find(req.Target).Text())

	resp := &model.Response{}
	resp.TaskID = req.TaskID
	resp.Body = body
	resp.Error = err

	return resp, nil
}
