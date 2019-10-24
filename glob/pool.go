package glob

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	go func() {
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
		<-stopChan
		for _, pool := range Pool.browsers {
			pool.Cancel()
		}
	}()
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
