package main

import (
	"gitlab.paradise-soft.com.tw/backend/legion/api"
	"gitlab.paradise-soft.com.tw/backend/legion/glob"
	"gitlab.paradise-soft.com.tw/backend/legion/service"
)

func main() {
	glob.Init()
	service.Init()
	api.Init()

	if err := api.Router.Run(glob.Config.WWW.Addr); err != nil {
		panic(err)
	}
}
