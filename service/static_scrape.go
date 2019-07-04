package service

import (
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

func StaticScrape(req model.Request) (interface{}, error) {
	doc, err := glob.GetAndConvertToDocument(req.URL)
	if err != nil {
		return nil, err
	}

	result := doc.Find(req.Target).Text()

	return result, nil
}
