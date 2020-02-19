package glob

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

var (
	port int = 9222
)

type Browser struct {
	Context   context.Context
	Cancel    context.CancelFunc
	Tabs      Tabs
	DebugPort int
	Options   []chromedp.ExecAllocatorOption
}

func NewBrowser() (*Browser, error) {
	b := &Browser{
		Tabs:      make(Tabs),
		DebugPort: port,
		Options: []chromedp.ExecAllocatorOption{
			chromedp.Flag("remote-debugging-port", fmt.Sprintf("%d", port)),
			chromedp.Flag("remote-debugging-address", "0.0.0.0"),
		},
	}

	port++
	b.Options = append(b.Options, browserOptions...)
	ctx, _ := chromedp.NewExecAllocator(context.Background(), b.Options...)

	b.Context, b.Cancel = chromedp.NewContext(ctx)

	if err := chromedp.Run(b.Context, chromedp.Navigate("about:blank")); err != nil {
		return nil, xerrors.Errorf("create browser error: %w", err)
	}
	return b, nil
}

func (b *Browser) NewTab(timeout time.Duration) (*Tab, error) {
	if err := b.Context.Err(); err != nil {
		ctx, _ := chromedp.NewExecAllocator(context.Background(), b.Options...)

		b.Context, b.Cancel = chromedp.NewContext(ctx)

		if err := chromedp.Run(b.Context, chromedp.Navigate("about:blank")); err != nil {
			return nil, xerrors.Errorf("recreate browser error: %w", err)
		}
	}
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
		return nil, xerrors.Errorf("create tab error: %w", err)
	}
	b.Tabs[uid] = tab
	return tab, nil
}
