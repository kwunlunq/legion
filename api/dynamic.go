package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
	"gitlab.paradise-soft.com.tw/dwh/legion/service"
)

func dynamicScrape(ctx *gin.Context) {
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

	resp, err := service.DynamicScrape(req)
	if err != nil {
		response(ctx, resp, http.StatusInternalServerError, -1, glob.ScrapeFailed, err)
		return
	}
	response(ctx, resp, http.StatusOK, 1, glob.ScrapeSuccess, nil)
}
