package glob

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
)

var (
	port int = 9222
)

type Browser struct {
	Context context.Context
	Cancel  context.CancelFunc
	Tabs    Tabs
}

func NewBrowser() (*Browser, error) {
	b := &Browser{
		Tabs: make(Tabs),
	}
	tmpBrowserOptions := []chromedp.ExecAllocatorOption{
		chromedp.Flag("remote-debugging-port", fmt.Sprintf("%d", port)),
	}
	port++
	tmpBrowserOptions = append(tmpBrowserOptions, browserOptions...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(),
		tmpBrowserOptions...)

	b.Context, b.Cancel = chromedp.NewContext(ctx)

	if err := chromedp.Run(b.Context, chromedp.Navigate("about:blank")); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Browser) NewTab(timeout time.Duration) (*Tab, error) {
	tab := &Tab{}
	uid := uuid.New().String()
	tab.ID = uid
	if timeout == 0 {
		timeout = Config.Chrome.Timeout
	}
	tab.orgContext, tab.orgCancel = context.WithTimeout(b.Context, timeout)
	tab.Context, tab.cancel = chromedp.NewContext(tab.orgContext)

	// tab.orgContext, tab.orgCancel = chromedp.NewContext(b.Context)
	// tab.Context, tab.cancel = context.WithTimeout(tab.orgContext, Config.Chrome.Timeout)
	if err := chromedp.Run(tab.Context, chromedp.Navigate("about:blank")); err != nil {
		return nil, err
	}
	b.Tabs[uid] = tab
	return tab, nil
}
