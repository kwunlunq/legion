package api

import (
	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/glob/dispatcher"
)

var Router *gin.Engine

func Init() {
	Router = gin.Default()
	InitAPIS()
	InitSubscribe()
}

func InitAPIS() {
	// /v1/apis
	apis := Router.Group(glob.Config.API.Version + "/apis")
	{
		apis.GET("/health", getHealth)

		// /v1/apis/dynamic
		dynamic := apis.Group("/dynamic")
		{
			dynamic.POST("/scrape", dynamicScrapeAPI)
			dynamic.GET("/cache", getDynamicCache)
		}

		// /v1/apis/static
		static := apis.Group("/static")
		{
			static.POST("/scrape", staticScrapeAPI)
			static.GET("/cache", getStaticCache)
		}
	}
}

func InitSubscribe() {
	dispatcher.Subscribe(glob.Config.Dispatcher.DynamicTopic,
		dynamicScrape,
		dispatcher.ConsumerOmitOldMsg(),
		dispatcher.ConsumerSetAsyncNum(glob.Config.Dispatcher.DynamicAsyncNum),
	)

	dispatcher.Subscribe(glob.Config.Dispatcher.StaticTopic,
		staticScrape,
		dispatcher.ConsumerOmitOldMsg(),
		dispatcher.ConsumerSetAsyncNum(glob.Config.Dispatcher.StaticAsyncNum),
	)
}
