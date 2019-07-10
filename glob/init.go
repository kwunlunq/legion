package glob

import (
	"time"

	"github.com/spf13/viper"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

var Config struct {
	Chrome struct {
		Headless      bool `mapstructure:"headless"`
		MaxBrowsers   int  `mapstructure:"max_browsers"`
		MaxTabs       int  `mapstructure:"max_tabs"`
		MaxRetryCount int  `mapstructure:"max_retry_count"`
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

func Init() {
	loadConfig()
	initTracer()
	initBrowserOptions()
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
