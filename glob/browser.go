package glob

import (
	"context"

	"github.com/chromedp/chromedp"
)

type Browser struct {
	Context context.Context
	Cancel  context.CancelFunc
	Tabs    Tabs
}

func NewBrowser() (*Browser, error) {
	b := &Browser{}

	ctx, _ := chromedp.NewExecAllocator(context.Background(), browserOptions...)

	b.Context, b.Cancel = chromedp.NewContext(ctx)
	if err := chromedp.Run(b.Context, chromedp.Navigate("about:blank")); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Browser) NewTab() (*Tab, error) {
	tab := &Tab{}
	tab.orgContext, tab.orgCancel = chromedp.NewContext(b.Context)
	tab.Context, tab.cancel = context.WithTimeout(tab.orgContext, Config.Chrome.Timeout)
	if err := chromedp.Run(tab.Context, chromedp.Navigate("about:blank")); err != nil {
		return nil, err
	}
	b.Tabs = append(b.Tabs, tab)
	return tab, nil
}
