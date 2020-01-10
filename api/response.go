package api

import (
	"fmt"
	"net/http"
	"time"

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

func logFormat(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	if param.StatusCode != http.StatusOK {
		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	} else {
		return ""
	}

}
