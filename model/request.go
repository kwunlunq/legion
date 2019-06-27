package model

type Request struct {
	TaskID    string  `json:"task_id"`
	RespTopic string  `json:"resp_topic"`
	URL       string  `json:"url"`
	Steps     []*Step `json:"steps"`
}
