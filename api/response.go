package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

type Response struct {
	Status  int         `json:"status"` // "success": 1, "fail": -1
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func response(ctx *gin.Context, item interface{}, success int, message string, err error) {
	resp := Response{
		Status:  success,
		Data:    item,
		Message: message,
	}

	if err != nil {
		tracer.Error("apis", err.Error())
		resp.Message = err.Error()
	}

	ctx.JSON(http.StatusOK, resp)
}

func responseParamError(ctx *gin.Context, err error) {
	resp := Response{
		Status: -1,
		// Data:    nil,
		Message: fmt.Sprintf("參數錯誤: %v", err.Error()),
	}
	ctx.JSON(http.StatusBadRequest, resp)
}
