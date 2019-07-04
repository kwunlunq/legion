package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

type Response struct {
	Status  int         `json:"status"` //1為成功、0為、-1為失敗
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func response(ctx *gin.Context, item interface{}, success int, message string, err error) {
	resp := Response{
		success,
		item,
		message,
	}

	if err != nil {
		tracer.Error("apis", err.Error())
	}

	ctx.JSON(http.StatusOK, resp)
}

func responseParamError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, Response{Message: "參數錯誤: " + err.Error(), Status: -1})
}
