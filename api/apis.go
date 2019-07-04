package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
)

var Router = gin.Default()

func Init() {
	apis := Router.Group(glob.Config.API.Version + "/apis")

	dynamic := apis.Group("/dynamic")
	dynamic.POST("/scrape", scrape)

	static := apis.Group("/static")
	static.POST("/scrape", scrape)
}
