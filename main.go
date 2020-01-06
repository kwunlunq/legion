package main

import (
	"log"
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
	go func() {
		if err := api.ListenTCP(8081); err != nil {
			panic(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-quit
	log.Println("Graceful shutdown start")
	glob.Pool.Close()
	log.Println("Graceful shutdown success")
}
