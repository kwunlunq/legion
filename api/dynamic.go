package api

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
	"gitlab.paradise-soft.com.tw/dwh/legion/service"
	"gitlab.paradise-soft.com.tw/glob/dispatcher"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

func dynamicScrape(data []byte) error {
	var err error
	var out []byte
	req := &model.Request{}

	if err = json.Unmarshal(data, req); err != nil {
		return err
	}

	if err = checkParams(req); err != nil {
		return err
	}

	if err = service.DynamicScrape(req); err != nil {
		return err
	}

	if out, err = json.Marshal(req); err != nil {
		return err
	}

	if err = dispatcher.Send(req.RespTopic, out, dispatcher.ProducerAddErrHandler(DispatcherErrHandler)); err != nil {
		return err
	}
	return nil
}

func DispatcherErrHandler(data []byte, err error) {
	if err != nil {
		tracer.Error("Dispatcher", err, data)
	}
}

func dynamicScrapeAPI(ctx *gin.Context) {
	req := &model.Request{}
	ctx.BindJSON(req)
	if err := checkParams(req); err != nil {
		responseParamError(ctx, err)
		return
	}

	if err := service.DynamicScrape(req); err != nil {
		response(ctx, req, -1, glob.ScrapeFailed, err)
		return
	}
	response(ctx, req, 1, glob.ScrapeSuccess, nil)
}

func getDynamicCache(ctx *gin.Context) {
	req := &model.CacheRequest{}

	ctx.BindQuery(req)

	if req.TaskID == "" {
		responseParamError(ctx, errors.New("task_id"))
		return
	}

	if err := service.GetDynamicCache(req); err != nil {
		response(ctx, req, -1, glob.ScrapeFailed, err)
		return
	}

	response(ctx, req, 1, glob.ScrapeSuccess, nil)
}
