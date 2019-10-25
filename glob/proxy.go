package glob

import (
	"errors"

	proxytool "gitlab.paradise-soft.com.tw/dwh/proxy/proxy"
)

var (
	ProxyDefaultAliveMinute = 30
)

var ErrCantGetProxy = errors.New("can't get proxy")

func initProxyService() {
	proxytool.InitProxyService("1000", "5f666d3c-8523-b9dbb86d-05c8-05df2de3", Config.ProxyService.Host)
}

func GetProxies(num int, countryCode []string) (proxies []string, err error) {
	collector := proxytool.NewCollector(
		proxytool.SetCountryCode(countryCode...),
		proxytool.SetAliveMin(ProxyDefaultAliveMinute),
		proxytool.SetNumber(num),
	)

	proxies, err = collector.GetProxys()
	if err != nil {
		return nil, err
	}

	return proxies, nil
}

func GetProxy(countryCode ...string) (proxy string, err error) {
	collector := proxytool.NewCollector(
		proxytool.SetCountryCode(countryCode...),
		proxytool.SetAliveMin(ProxyDefaultAliveMinute),
		proxytool.SetNumber(1),
	)

	var proxies []string
	proxies, err = collector.GetProxys()
	if err != nil {
		return "", err
	}

	if len(proxies) == 0 {
		return "", ErrCantGetProxy
	}

	return proxies[0], nil
}
