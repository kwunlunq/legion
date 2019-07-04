package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
)

func getHealth(ctx *gin.Context) {
	response(ctx, nil, http.StatusOK, 1, glob.ServiceSuccess, nil)
}
