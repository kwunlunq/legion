package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gitlab.paradise-soft.com.tw/dwh/legion/api"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/service"
)

func main() {
	glob.Init()
	service.Init()
	api.Init()
	go func() {
		if err := api.Router.Run(glob.Config.WWW.Addr); err != nil {
			panic(err)
		}
	}()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	glob.Pool.Close()
	fmt.Println("1234")
}
