package service

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
)

const (
	Click       = "click"
	DoubleClick = "double_cilck"
	SendKeys    = "send_keys"
	WaitReady   = "wait_ready"
	WaitVisible = "wait_visible"
)

func (this *LegionRequest) GetDynamicResponse() (LegionResp *LegionResponse) {
	var body []byte
	var ReqErr error
	body, ReqErr = this.doDynamic()

	LegionResp = &LegionResponse{}
	LegionResp.Req = this

	if ReqErr != nil {
		LegionResp.ErrorMessages = []string{ReqErr.Error()}
		return
	}

	LegionResp.Body = body
	return
}

// DynamicRequest
func (this *LegionRequest) toDynamicRequest() (dynamicReq *DynamicRequest, err error) {
	dynamicReq = &DynamicRequest{}
	dynamicReq.RawURL = this.RawURL
	dynamicReq.Target = this.Target
	dynamicReq.Steps = this.Steps
	return
}

func (this *LegionRequest) doDynamic() (body []byte, err error) {
	defer func() {
		if err != nil {
			body = nil
		}
	}()

	var dynamicReq *DynamicRequest
	dynamicReq, err = this.toDynamicRequest()
	if err != nil {
		return
	}

	// Todo: err is not handled correctly
	tab := glob.Pool.NewTab()
	defer func() {
		tab.Cancel()
		glob.Pool.RemoveTab(tab)
	}()

	body, err = dynamicReq.runTasks(tab.Context)
	if err != nil {
		return
	}

	if this.Charset != "" {
		body, err = glob.Decoder(body, this.Charset)
		if err != nil {
			return
		}
	}

	return
}

type DynamicRequest struct {
	RawURL string  `json:"rawURL"`
	Steps  []*Step `json:"steps"`
	Target string  `json:"target"`
}

type Step struct {
	Action string `json:"action"`
	Target string `json:"target"`
	Keys   string `json:"keys"`
}

func (req *DynamicRequest) runTasks(ctx context.Context) ([]byte, error) {
	tasks, err := req.makeTasks(req.Steps)
	if err != nil {
		return nil, err
	}

	if err = chromedp.Run(ctx, chromedp.Navigate(req.RawURL)); err != nil {
		err = fmt.Errorf(`%s while navigating to "%s"`, err.Error(), req.RawURL)
		return nil, err
	}

	doneSteps := 0
	for _, task := range tasks {
		if err = chromedp.Run(ctx, task); err != nil {
			err = fmt.Errorf(`%s while executing step[%d] "%s %s"`, err.Error(), doneSteps+1, req.Steps[doneSteps].Action, req.Steps[doneSteps].Target)
			return nil, err
		}
		doneSteps++
	}

	var result string
	if err = chromedp.Run(ctx, chromedp.OuterHTML(req.Target, &result)); err != nil {
		err = fmt.Errorf(`%s while retrieving outer html from "%s"`, err.Error(), req.Target)
		return nil, err
	}

	return []byte(result), nil
}

func (req *DynamicRequest) makeTasks(steps []*Step) (chromedp.Tasks, error) {
	var err error
	tasks := chromedp.Tasks{}

	for _, step := range steps {
		switch step.Action {
		case Click:
			tasks = append(tasks, chromedp.Click(step.Target))
		case DoubleClick:
			tasks = append(tasks, chromedp.DoubleClick(step.Target))
		case SendKeys:
			tasks = append(tasks, chromedp.SendKeys(step.Target, step.Keys))
		case WaitReady:
			tasks = append(tasks, chromedp.WaitReady(step.Target))
		case WaitVisible:
			tasks = append(tasks, chromedp.WaitVisible(step.Target))
		default:
			err = fmt.Errorf(`Unsupported step action "%s"`, step.Action)
			return nil, err
		}
	}

	return tasks, nil
}
