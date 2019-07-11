package api

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/glob/dispatcher"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
	"gitlab.paradise-soft.com.tw/dwh/legion/service"
)

func staticScrape(data []byte) error {
	var err error
	var out []byte
	req := &model.Request{}

	if err = json.Unmarshal(data, req); err != nil {
		return err
	}

	if err = checkParams(req); err != nil {
		return err
	}

	if err = service.StaticScrape(req); err != nil {
		return err
	}

	if out, err = json.Marshal(req); err != nil {
		return err
	}

	if err = dispatcher.Send(req.RespTopic, out); err != nil {
		return err
	}
	return nil
}

func staticScrapeAPI(ctx *gin.Context) {
	req := &model.Request{}

	ctx.BindJSON(req)

	if err := checkParams(req); err != nil {
		responseParamError(ctx, err)
		return
	}

	if err := service.StaticScrape(req); err != nil {
		response(ctx, req, -1, glob.ScrapeFailed, err)
		return
	}

	response(ctx, req, 1, glob.ScrapeSuccess, nil)
}

func getStaticCache(ctx *gin.Context) {
	req := &model.CacheRequest{}

	ctx.BindQuery(req)

	if req.TaskID == "" {
		responseParamError(ctx, errors.New("task_id"))
		return
	}

	if err := service.GetStaticCache(req); err != nil {
		response(ctx, req, -1, glob.ScrapeFailed, err)
		return
	}

	response(ctx, req, 1, glob.ScrapeSuccess, nil)
}
