package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
)

var Router = gin.Default()

func Init() {
	InitAPIS()
}

func InitAPIS() {
	apis := Router.Group(glob.Config.API.Version + "/apis")

	apis.GET("/health", getHealth)

	dynamic := apis.Group("/dynamic")
	dynamic.POST("/scrape", dynamicScrapeAPI)
	dynamic.GET("/cache", getDynamicCache)

	static := apis.Group("/static")
	static.POST("/scrape", staticScrapeAPI)
	static.GET("/cache", getStaticCache)
}
