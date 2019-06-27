package service

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

func Scrape(req model.Request) (model.Response, error) {
	ctx, cancel := glob.NewTabContext()
	defer cancel()
	body, err := runTasks(ctx, req)

	resp := model.Response{}
	resp.TaskID = req.TaskID
	resp.Body = body
	if err != nil {
		resp.Error = err.Error()
	}
	return resp, err
}

func runTasks(ctx context.Context, req model.Request) ([]byte, error) {
	resp := model.Response{}
	resp.TaskID = req.TaskID

	tasks, err := makeTasks(req.Steps)
	if err != nil {
		return nil, err
	}

	if err = chromedp.Run(ctx, chromedp.Navigate(req.URL)); err != nil {
		err = fmt.Errorf(`%s while navigating "%s"`, err.Error(), req.URL)
		return nil, err
	}

	var html string
	tasks = append(tasks)

	doneSteps := 0
	for _, task := range tasks {
		if err = chromedp.Run(ctx, task); err != nil {
			err = fmt.Errorf(`%s while executing step[%d] "%s %s"`, err.Error(), doneSteps+1, req.Steps[doneSteps].Action, req.Steps[doneSteps].Target)
			return nil, err
		}
		doneSteps++
	}

	if err = chromedp.Run(ctx, chromedp.OuterHTML(`html`, &html)); err != nil {
		err = fmt.Errorf("%s while retrieving html", err.Error())
		return nil, err
	}

	return []byte(html), err
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
