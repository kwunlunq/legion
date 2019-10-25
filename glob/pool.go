package glob

import (
	"sync"

	"gitlab.paradise-soft.com.tw/glob/tracer"
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
	for len(p.browsers) < p.maxBrowsers {
		if retryCount > p.maxRetryCount {
			return
		}

		b, err := NewBrowser()
		if err != nil {
			retryCount++
			tracer.Error("NewBrowser", err)
			continue
		}
		p.browsers = append(p.browsers, b)
	}
}

func (p *pool) NewTab() *Tab {
	p.Lock()
	defer p.Unlock()

	var tab *Tab
	for _, b := range p.browsers {
		if len(b.Tabs) < p.maxTabs {
			tab, _ = b.NewTab()
			break
		}
	}

	return tab
}

func (p *pool) RemoveTab(tab *Tab) {
	p.Lock()
	defer p.Unlock()

	for _, b := range p.browsers {
		for i, t := range b.Tabs {
			if t == tab {
				b.Tabs = append(b.Tabs[:i], b.Tabs[i+1:]...)
				return
			}
		}
	}
}

func (p *pool) Close() {
	p.Lock()
	defer p.Unlock()

	for _, p := range p.browsers {
		p.Cancel()
	}
}
