package glob

import (
	"sync"

	"github.com/chromedp/chromedp"
)

type pool struct {
	browsers    []*Browser
	maxBrowsers int
	maxTabs     int
}

func NewPool(browsersNum, tabsNum int) *pool {
	p := &pool{}
	p.Fill(browsersNum, tabsNum)
	return p
}

func (p *pool) Fill(browsersNum, tabsNum int) {
	for len(p.browsers) < browsersNum {
		b := NewBrowser()
		p.browsers = append(p.browsers, b)
		for j := 0; j < tabsNum; j++ {
			tab, _ := b.NewTab()
		}
	}
}

func (p *pool) GetBrowser() *Browser {
	browser := &Browser{}
	for b := range p.browserToTabs {
		targets, _ := chromedp.Targets(b.Context)
		if len(targets) < p.maxTabs+5 {
			browser = b
		}
	}
	return browser
}

func (p *pool) GetTab() *Tab {
	tab := &Tab{}
	mutex := sync.Mutex{}
	mutex.Lock()
	for b := range p.browserToTabs {
		if len(p.browserToTabs[b]) < p.maxTabs {
			tab, _ = NewTab(b.Context)
			p.browserToTabs[b] = append(p.browserToTabs[b], tab)
			break
		}
	}
	mutex.Unlock()
	return tab
}

func (p *pool) RemoveTab(*Tab) {
	tab := &Tab{}
	mutex := sync.Mutex{}
	mutex.Lock()
	for b := range p.browserToTabs {
		if len(p.browserToTabs[b]) < p.maxTabs {
			tab, _ = NewTab(b.Context)
			break
		}
	}
	mutex.Unlock()
	return tab
}
