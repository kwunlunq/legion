package glob

import (
	"context"
)

type Tabs []*Tab

type Tab struct {
	Context context.Context
	Cancel  context.CancelFunc
	Browser *Browser
}
