package api

import (
	"errors"

	"gitlab.paradise-soft.com.tw/dwh/legion/model"
)

func checkParams(req *model.Request) error {
	if req.TaskID == "" {
		return errors.New("task_id")
	}
	if req.URL == "" {
		return errors.New("url")
	}
	if req.RespTopic == "" {
		return errors.New("resp_topic")
	}
	if req.Target == "" {
		return errors.New("target")
	}
	return nil
}
