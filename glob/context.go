package glob

import (
	"context"

	"github.com/chromedp/chromedp"
)

func NewTabContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(DefaultBrowserCTX, Config.API.Timeout)
	ctx, _ = chromedp.NewContext(ctx)
	return ctx, cancel
}

func NewBrowserContext() (context.Context, context.CancelFunc) {
	opts := append(DefaultExecAllocatorOptions[:], chromedp.ExecPath(Config.Chrome.Path))
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = chromedp.NewContext(ctx)
	return ctx, cancel
}
