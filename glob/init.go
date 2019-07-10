package glob

import (
	"time"

	"github.com/spf13/viper"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

var Config struct {
	Chrome struct {
		Path          string `mapstructure:"path"`
		MaxBrowsers   int    `mapstructure:"max_browsers"`
		MaxTabs       int    `mapstructure:"max_tabs"`
		MaxRetryCount int    `mapstructure:"max_retry_count"`
	} `mapstructure:"chrome"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
	WWW struct {
		Addr string `mapstructure:"addr"`
		Host string `mapstructure:"host"`
	} `mapstructure:"www"`
	API struct {
		Version string        `mapstructure:"version"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"api"`
	CPU struct {
		Limit int `mapstructure:"limit"`
	} `mapstructure:"cpu"`
}

var (
	Pool  *pool
	Cache *cache
)

func Init() {
	loadConfig()
	initTracer()
	initPool()
	initCache()
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

func initPool() {
	Pool = NewPool(Config.Chrome.MaxBrowsers, Config.Chrome.MaxTabs)
}

func initCache() {
	Cache = NewCache()
}
