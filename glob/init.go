package glob

import (
	"net/url"
	"strings"
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
	GoRequest struct {
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"goRequest"`
	Log struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
	WWW struct {
		Addr         string `mapstructure:"addr"`
		Host         string `mapstructure:"host"`
		InternalHost string `mapstructure:"internalHost"`
		ExternalHost string `mapstructure:"externalHost"`
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
	ProxyService struct {
		Host string `mapstructure:"host"`
	} `mapstructure:"proxyService"`
}

func Init() {
	loadConfig()
	initTracer()
	// initBrowserOptions()
	// initPool()
	initRespCache()
	initDispatcher()
	initProxyService()
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

	Config.WWW.InternalHost = strings.TrimRight(Config.WWW.InternalHost, "/")
	Config.WWW.ExternalHost = strings.TrimRight(Config.WWW.ExternalHost, "/")
	Config.Log.Level = strings.ToLower(Config.Log.Level)
}

func initTracer() {
	tracer.SetLevelWithName(Config.Log.Level)
}

func GetInternalHostURL() (u *url.URL) {
	var err error
	u, err = url.Parse(Config.WWW.InternalHost)
	if err != nil {
		panic(err)
	}
	return u
}

func GetExternalHostURL() (u *url.URL) {
	var err error
	u, err = url.Parse(Config.WWW.ExternalHost)
	if err != nil {
		panic(err)
	}
	return u
}
