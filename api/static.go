package api

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
	"gitlab.paradise-soft.com.tw/dwh/legion/service"
)

func staticScrape(ctx *gin.Context) {
	req := model.Request{}
	ctx.BindJSON(&req)
	if req.TaskID == "" {
		responseParamError(ctx, errors.New("task_id"))
		return
	}
	if req.URL == "" {
		responseParamError(ctx, errors.New("url"))
		return
	}
	if req.RespTopic == "" {
		responseParamError(ctx, errors.New("resp_topic"))
		return
	}
	if req.Target == "" {
		responseParamError(ctx, errors.New("target"))
		return
	}

	resp, err := service.StaticScrape(req)
	if err != nil {
		response(ctx, resp, -1, glob.ScrapeFailed, err)
		return
	}
	response(ctx, resp, 1, glob.ScrapeSuccess, nil)
}
