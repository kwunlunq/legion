package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/glob/dispatcher"
	"gitlab.paradise-soft.com.tw/glob/tracer"
	"golang.org/x/sync/errgroup"
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
	go func() {
		eg, _ := errgroup.WithContext(context.Background())
		eg.Go(func() (err error) {
			err = dispatcher.Subscribe(glob.Config.Dispatcher.DynamicTopic,
				dynamicScrape,
				dispatcher.ConsumerOmitOldMsg(),
				dispatcher.ConsumerSetAsyncNum(glob.Config.Dispatcher.DynamicAsyncNum),
			)
			// if err != nil {
			tracer.Errorf("dispatcher", "dynamic scrape: %v", err)
			// }
			return
		})

		eg.Go(func() (err error) {
			err = dispatcher.Subscribe(glob.Config.Dispatcher.StaticTopic,
				staticScrape,
				dispatcher.ConsumerOmitOldMsg(),
				dispatcher.ConsumerSetAsyncNum(glob.Config.Dispatcher.StaticAsyncNum),
			)
			// if err != nil {
			tracer.Errorf("dispatcher", "static scrape: %v", err)
			// }
			return
		})

		err := eg.Wait()
		tracer.Errorf("dispatcher", "scrape error: %v", err)
	}()
}
