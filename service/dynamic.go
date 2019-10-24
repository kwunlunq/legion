package service

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	sdk "gitlab.paradise-soft.com.tw/glob/legion-sdk"
)

func (r *LegionRequest) GetDynamicResult() (legionResult *LegionResult) {
	// var resp *http.Response
	var body []byte
	var err error
	body, err = r.doDynamic()

	legionResult = &LegionResult{}
	legionResult.Request = (*sdk.LegionRequest)(r)
	if err != nil {
		legionResult.ErrorMessage = err.Error()
		return
	}

	legionResp := &LegionResponse{}
	legionResp.Body = body
	legionResult.Response = (*sdk.LegionResponse)(legionResp)
	return
}

// DynamicRequest
func (r *LegionRequest) toDynamicRequest() (dynamicReq *DynamicRequest, err error) {
	dynamicReq = &DynamicRequest{}
	dynamicReq.RawURL = r.RawURL
	dynamicReq.Target = r.Target
	dynamicReq.Steps = r.Steps
	return
}

func (r *LegionRequest) doDynamic() (body []byte, err error) {
	defer func() {
		if err != nil {
			body = nil
		}
	}()

	var dynamicReq *DynamicRequest
	dynamicReq, err = r.toDynamicRequest()
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

	if r.Charset != "" {
		body, err = glob.Decoder(body, r.Charset)
		if err != nil {
			return
		}
	}

	return
}

type DynamicRequest struct {
	RawURL string      `json:"rawURL"`
	Steps  []*sdk.Step `json:"steps"`
	Target string      `json:"target"`
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

func (req *DynamicRequest) makeTasks(steps []*sdk.Step) (chromedp.Tasks, error) {
	var err error
	tasks := chromedp.Tasks{}

	for _, step := range steps {
		switch step.Action {
		case sdk.Click:
			tasks = append(tasks, chromedp.Click(step.Target))
		case sdk.DoubleClick:
			tasks = append(tasks, chromedp.DoubleClick(step.Target))
		case sdk.SendKeys:
			tasks = append(tasks, chromedp.SendKeys(step.Target, step.Keys))
		case sdk.WaitReady:
			tasks = append(tasks, chromedp.WaitReady(step.Target))
		case sdk.WaitVisible:
			tasks = append(tasks, chromedp.WaitVisible(step.Target))
		default:
			err = fmt.Errorf(`Unsupported step action "%s"`, step.Action)
			return nil, err
		}
	}

	return tasks, nil
}
