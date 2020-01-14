package glob

import (
	"context"
)

type Tabs map[string]*Tab

type Tab struct {
	ID         string
	Context    context.Context
	cancel     context.CancelFunc
	orgContext context.Context
	orgCancel  context.CancelFunc
	// Browser    *Browser
}

func (t *Tab) Cancel() {
	if t.Context.Err() == nil {
		t.cancel()
	}
	if t.orgContext.Err() == nil {
		t.orgCancel()
	}
}
