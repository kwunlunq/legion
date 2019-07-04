package glob

import (
	"context"

	"github.com/chromedp/chromedp"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

type Browser struct {
	Context context.Context
	Cancel  context.CancelFunc
	Tabs    Tabs
}

func NewBrowser() *Browser {
	b := &Browser{}
	opts := append(DefaultExecAllocatorOptions[:], chromedp.ExecPath(Config.Chrome.Path))
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	b.Context, b.Cancel = chromedp.NewContext(ctx)
	if err := chromedp.Run(ctx, chromedp.Navigate("about:blank")); err != nil {
		tracer.Error("NewBrowser", err)
	}
	return b
}

func (b *Browser) NewTab() (*Tab, error) {
	tab := &Tab{}
	tab.Context, tab.Cancel = chromedp.NewContext(b.Context)
	if err := chromedp.Run(tab.Context, chromedp.Navigate("about:blank")); err != nil {
		return nil, err
	}
	b.Tabs = append(b.Tabs, tab)
	return tab, nil
}
