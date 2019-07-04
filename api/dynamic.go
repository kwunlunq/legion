package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
	"gitlab.paradise-soft.com.tw/dwh/legion/service"
	"gitlab.paradise-soft.com.tw/glob/common/codebook"
)

func scrape(ctx *gin.Context) {
	req := model.Request{}
	ctx.BindJSON(&req)
	if req.TaskID == "" {
		err := glob.StatusDetail(codebook.Status_Arguments_Missing, "task_id")
		response(ctx, nil, http.StatusBadRequest, -1, err.Error(), nil)
		return
	}
	if req.URL == "" {
		err := glob.StatusDetail(codebook.Status_Arguments_Missing, "url")
		response(ctx, nil, http.StatusBadRequest, -1, err.Error(), nil)
		return
	}
	if req.RespTopic == "" {
		err := glob.StatusDetail(codebook.Status_Arguments_Missing, "resp_topic")
		response(ctx, nil, http.StatusBadRequest, -1, err.Error(), nil)
		return
	}

	resp, err := service.Scrape(req)
	if err != nil {
		response(ctx, resp, http.StatusInternalServerError, -1, glob.ScrapeFailed, err)
		return
	}
	response(ctx, resp, http.StatusOK, 1, glob.ScrapeSuccess, nil)
}
