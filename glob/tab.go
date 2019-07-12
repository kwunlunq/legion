package glob

import (
	"context"
)

type Tabs []*Tab

type Tab struct {
	Context    context.Context
	cancel     context.CancelFunc
	orgContext context.Context
	orgCancel  context.CancelFunc
	Browser    *Browser
}

func (t *Tab) Cancel() {
	if t.Context.Err() == nil {
		t.cancel()
	}
	if t.orgContext.Err() == nil {
		t.orgCancel()
	}
}
