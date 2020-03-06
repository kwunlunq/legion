package glob

import (
	"sync"
	"time"

	"gitlab.paradise-soft.com.tw/glob/tracer"
	"golang.org/x/xerrors"
)

var (
	Pool *pool
)

type pool struct {
	sync.RWMutex
	browsers      []*Browser
	maxBrowsers   int
	maxTabs       int
	maxRetryCount int
}

func initPool() {
	Pool = NewPool(Config.Chrome.MaxBrowsers, Config.Chrome.MaxTabs)

}

func NewPool(maxBrowsers, maxTabs int) *pool {
	maxRetryCount := Config.Chrome.MaxRetryCount
	p := &pool{
		maxBrowsers:   maxBrowsers,
		maxTabs:       maxTabs,
		maxRetryCount: maxRetryCount,
	}
	p.Fill()
	return p
}

func (p *pool) Fill() {
	// Todo: Sometimes the program gets stuck in this function
	p.Lock()
	defer p.Unlock()

	retryCount := 0
	port := 9222

	for len(p.browsers) < p.maxBrowsers {
		if retryCount > p.maxRetryCount {
			return
		}
		var browser *Browser
		var err error
		if len(p.browsers) == 0 {
			browser, err = NewBrowser(SetUseProxy(false), SetRemoteDebugging(port))
		} else {
			browser, err = NewBrowser(SetUseProxy(true), SetRemoteDebugging(port))
		}
		if err != nil {
			retryCount++
			tracer.Error("NewBrowser", err)
			continue
		}
		p.browsers = append(p.browsers, browser)
		port++
	}
}

func (p *pool) NewTab(isUseProxy bool, timeout time.Duration) *Tab {
	p.Lock()
	defer p.Unlock()

	var tab *Tab
	for i, b := range p.browsers {
		if isUseProxy == b.IsUseProxy {
			if len(b.Tabs) < p.maxTabs {
				var err error
				tab, err = b.NewTab(timeout)
				if err != nil {
					tracer.Errorf("tab", "create tab error: %s", err)
					if xerrors.Is(err, ErrorBrowserContext) || xerrors.Is(err, ErrorTabContext) {
						b.Cancel()
						p.browsers[i], err = NewBrowser(SetUseProxy(b.IsUseProxy), SetRemoteDebugging(b.DebugPort))
					}
					continue
				}
				return tab
			}
		}

	}

	return tab
}

func (p *pool) RemoveTab(tab *Tab) {
	p.Lock()
	defer p.Unlock()

	for _, b := range p.browsers {
		if _, ok := b.Tabs[tab.ID]; ok {
			delete(b.Tabs, tab.ID)
		}
		// for i, t := range b.Tabs {
		// 	if t == tab {
		// 		b.Tabs = append(b.Tabs[:i], b.Tabs[i+1:]...)
		// 		return
		// 	}
		// }
	}
}

func (p *pool) Close() {
	p.Lock()
	defer p.Unlock()

	for _, p := range p.browsers {
		p.Cancel()
	}
}

func (p *pool) GetBrowsersInfo() []*Browser {
	return p.browsers
}
