package glob

import (
	"time"

	"github.com/spf13/viper"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

var Config struct {
	Chrome struct {
		Headless      bool          `mapstructure:"headless"`
		MaxBrowsers   int           `mapstructure:"max_browsers"`
		MaxTabs       int           `mapstructure:"max_tabs"`
		MaxRetryCount int           `mapstructure:"max_retry_count"`
		Timeout       time.Duration `mapstructure:"timeout"`
	} `mapstructure:"chrome"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
	WWW struct {
		Addr string `mapstructure:"addr"`
		Host string `mapstructure:"host"`
	} `mapstructure:"www"`
	API struct {
		Version string `mapstructure:"version"`
	} `mapstructure:"api"`
	CPU struct {
		Limit int `mapstructure:"limit"`
	} `mapstructure:"cpu"`
	Dispatcher struct {
		Brokers         []string `mapstructure:"brokers"`
		GroupID         string   `mapstructure:"group_id"`
		DynamicTopic    string   `mapstructure:"dynamic_topic"`
		StaticTopic     string   `mapstructure:"static_topic"`
		DynamicAsyncNum int      `mapstructure:"dynamic_async_num"`
		StaticAsyncNum  int      `mapstructure:"static_async_num"`
	} `mapstructure:"dispatcher"`
}

func Init() {
	loadConfig()
	initTracer()
	initBrowserOptions()
	initPool()
	initCache()
	initDispatcher()
}

func loadConfig() {
	viper.SetConfigFile("app.conf")
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}

func initTracer() {
	tracer.SetLevelWithName(Config.Log.Level)
}
