package service

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

func DynamicScrape(req model.Request) (*model.Response, error) {
	tab := glob.Pool.NewTab()
	defer func() {
		tab.Cancel()
		glob.Pool.RemoveTab(tab)
	}()

	body, err := runTasks(tab.Context, req)

	if req.Charset != "" {
		body, err = glob.Decoder(body, req.Charset)
	}

	resp := &model.Response{}
	resp.TaskID = req.TaskID
	resp.Body = body
	resp.Error = err

	if err := glob.Cache.SetDynamicCache(req.TaskID, body); err != nil {
		return nil, err
	}

	return resp, err
}

func GetDynamicCache(req model.CacheRequest) (*model.CacheResponse, error) {
	resp := &model.CacheResponse{}

	resp.TaskID = req.TaskID

	value, err := glob.Cache.GetDynamicCache(req.TaskID)
	if err != nil {
		return nil, err
	}

	resp.Content = string(value)

	return resp, nil
}

func runTasks(ctx context.Context, req model.Request) ([]byte, error) {
	tasks, err := makeTasks(req.Steps)
	if err != nil {
		return nil, err
	}

	if err = chromedp.Run(ctx, chromedp.Navigate(req.URL)); err != nil {
		err = fmt.Errorf(`%s while navigating to "%s"`, err.Error(), req.URL)
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

func makeTasks(steps []*model.Step) (chromedp.Tasks, error) {
	var err error
	tasks := chromedp.Tasks{}

	for _, step := range steps {
		switch step.Action {
		case model.Click:
			tasks = append(tasks, chromedp.Click(step.Target))
		case model.DoubleClick:
			tasks = append(tasks, chromedp.DoubleClick(step.Target))
		case model.SendKeys:
			tasks = append(tasks, chromedp.SendKeys(step.Target, step.Keys))
		case model.WaitReady:
			tasks = append(tasks, chromedp.WaitReady(step.Target))
		case model.WaitVisible:
			tasks = append(tasks, chromedp.WaitVisible(step.Target))
		default:
			err = fmt.Errorf(`Unsupported step action "%s"`, step.Action)
			return nil, err
		}
	}

	return tasks, nil
}