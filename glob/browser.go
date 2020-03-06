package glob

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

var (
	ErrorBrowserContext = errors.New("browser ctx error")
	ErrorTabContext     = errors.New("tab ctx error")
)

type Browser struct {
	Context    context.Context
	Cancel     context.CancelFunc
	Tabs       Tabs
	DebugPort  int
	Options    []chromedp.ExecAllocatorOption
	IsUseProxy bool
}

func NewBrowser(opts ...BrowserOption) (*Browser, error) {
	b := &Browser{
		Tabs:    make(Tabs),
		Options: []chromedp.ExecAllocatorOption{},
	}

	b.Options = append(b.Options, browserOptions...)

	for _, opt := range opts {
		err := opt(b)
		if err != nil {
			return nil, err
		}
	}

	ctx, _ := chromedp.NewExecAllocator(context.Background(), b.Options...)

	b.Context, b.Cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))

	if err := chromedp.Run(b.Context, chromedp.Navigate("about:blank")); err != nil {
		return nil, xerrors.Errorf("create browser error: %w", err)
	}
	return b, nil
}

func (b *Browser) NewTab(timeout time.Duration) (*Tab, error) {
	if err := b.Context.Err(); err != nil {
		err = xerrors.Errorf("err: %w, %s", ErrorBrowserContext, err)
		return nil, err
		// browser, err = NewBrowser(SetUseProxy(false), SetRemoteDebugging(port))
		//
		// ctx, _ := chromedp.NewExecAllocator(context.Background(), b.Options...)
		//
		// b.Context, b.Cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
		//
		// if err := chromedp.Run(b.Context, chromedp.Navigate("about:blank")); err != nil {
		// 	return nil, xerrors.Errorf("recreate browser error: %w", err)
		// }
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
		err = xerrors.Errorf("err: %w, %s", ErrorTabContext, err)
		return nil, err
	}
	b.Tabs[uid] = tab
	return tab, nil
}

type BrowserOption func(c *Browser) error

func SetUseProxy(isUse bool) BrowserOption {
	return func(c *Browser) error {
		c.IsUseProxy = isUse
		if isUse {
			c.Options = append(c.Options, chromedp.ProxyServer("http://127.0.0.1:8081"))
		}
		return nil
	}
}

func SetRemoteDebugging(port int) BrowserOption {
	return func(c *Browser) error {
		c.DebugPort = port
		c.Options = append(c.Options, chromedp.Flag("remote-debugging-port", fmt.Sprintf("%d", port)))
		c.Options = append(c.Options, chromedp.Flag("remote-debugging-address", "0.0.0.0"))
		return nil
	}
}
