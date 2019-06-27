package glob

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/spf13/viper"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

var Config struct {
	Chrome struct {
		Path string `toml:"path"`
	} `toml:"chrome"`
	Log struct {
		Level string `toml:"level"`
	} `toml:"log"`
	WWW struct {
		Addr string `toml:"addr"`
		Host string `toml:"host"`
	} `toml:"www"`
	API struct {
		Version string        `toml:"version"`
		Timeout time.Duration `toml:"timeout"`
	} `toml:"api"`
	CrawlerSetting struct {
		IsAsync            int `toml:"is_async"`
		DefaultParallelism int `toml:"default_parallelism"`
		RequestTimeout     int `toml:"request_timeout"`
		CPULimit           int `toml:"cpu_limit"`
	} `toml:"crawler_setting"`
	Dispatcher struct {
		Brokers []string `toml:"brokers"`
		GroupID string   `toml:"group_id"`
	} `toml:"dispatcher"`
	ProxyService struct {
		Host string `toml:"host"`
	} `toml:"proxyService"`
}

var (
	DefaultBrowserCTX context.Context
)

func Init() {
	loadConfig()
	initTracer()
	InitContext()
}

func loadConfig() {
	viper.SetConfigFile("app.conf")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.Unmarshal(&Config)
}

func initTracer() {
	tracer.SetLevelWithName(Config.Log.Level)
}

func InitContext() {
	DefaultBrowserCTX, _ = NewBrowserContext()
	if err := chromedp.Run(DefaultBrowserCTX, chromedp.Navigate("http://www.google.com")); err != nil {
		panic(err)
	}
}
