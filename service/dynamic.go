package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/glob/helper"
	sdk "gitlab.paradise-soft.com.tw/glob/legion-sdk"
)

func (r *LegionRequest) GetDynamicResult() (legionResult *LegionResult) {
	// var resp *http.Response
	var response *network.Response
	var body []byte
	var cookies []*network.Cookie
	var err error
	response, body, cookies, err = r.doDynamic()

	legionResult = &LegionResult{}
	legionResult.Request = (*sdk.LegionRequest)(r)
	if err != nil {
		legionResult.ErrorMessage = err.Error()
		return
	}

	legionResp := &LegionResponse{}
	legionResp.FinishedAt = helper.Now(8)
	legionResp.StatusCode = int(response.Status)
	legionResp.Header = make(map[string]string, len(response.Headers))
	for key, val := range response.Headers {
		legionResp.Header[key] = fmt.Sprintf("%v", val)
	}
	legionResp.RequestHeader = make(map[string]string, len(response.RequestHeaders))
	for key, val := range response.RequestHeaders {
		legionResp.RequestHeader[key] = fmt.Sprintf("%v", val)
	}
	legionResp.Cookies = cookies
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

func (r *LegionRequest) doDynamic() (response *network.Response, body []byte, cookies []*network.Cookie, err error) {
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
	var tab *glob.Tab
	tab = glob.Pool.NewTab(len(r.ProxyLocations) > 0, r.Timeout)
	for tab == nil {
		tab = glob.Pool.NewTab(len(r.ProxyLocations) > 0, r.Timeout)
		time.Sleep(1 * time.Second)
	}
	defer func() {
		tab.Cancel()
		glob.Pool.RemoveTab(tab)
	}()

	response, body, cookies, err = dynamicReq.runTasks(tab.Context)
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

func (req *DynamicRequest) runTasks(ctx context.Context) (response *network.Response, body []byte, cookies []*network.Cookie, err error) {
	tasks, err := req.makeTasks(req.Steps)
	if err != nil {
		return nil, nil, nil, err
	}
	response = &network.Response{}

	if err = chromedp.Run(ctx, chromeTask(ctx, req.RawURL, response)); err != nil {
		err = fmt.Errorf(`%s while navigating to "%s"`, err.Error(), req.RawURL)
		return nil, nil, nil, err
	}

	doneSteps := 0
	for _, task := range tasks {
		if err = chromedp.Run(ctx, task); err != nil {
			err = fmt.Errorf(`%s while executing step[%d] "%s %s"`, err.Error(), doneSteps+1, req.Steps[doneSteps].Action, req.Steps[doneSteps].Target)
			return nil, nil, nil, err
		}
		doneSteps++
	}

	var result string
	if err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			data, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			cookies = data
			return nil
		}),
		chromedp.OuterHTML(req.Target, &result)); err != nil {
		err = fmt.Errorf(`%s while retrieving outer html from "%s"`, err.Error(), req.Target)
		return nil, nil, nil, err
	}

	return response, []byte(result), cookies, nil
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
		case sdk.Sleep:
			d, err := time.ParseDuration(step.Target)
			if err != nil {
				continue
			}
			tasks = append(tasks, chromedp.Sleep(d))
		case sdk.Reload:
			tasks = append(tasks, chromedp.Reload())

		default:
			err = fmt.Errorf(`Unsupported step action "%s"`, step.Action)
			return nil, err
		}
	}

	return tasks, nil
}

func chromeTask(chromeContext context.Context, url string, response *network.Response) chromedp.Tasks {
	chromedp.ListenTarget(chromeContext, func(event interface{}) {
		switch responseReceivedEvent := event.(type) {
		case *network.EventResponseReceived:
			if responseReceivedEvent.Response.URL == url {
				*response = *(responseReceivedEvent.Response)
			}
		}
	})

	return chromedp.Tasks{
		network.Enable(),
		chromedp.Navigate(url),
	}
}
