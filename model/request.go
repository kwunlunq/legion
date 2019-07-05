package model

type Request struct {
	TaskID    string  `json:"task_id"`
	RespTopic string  `json:"resp_topic"`
	URL       string  `json:"url"`
	Target    string  `json:"target"`
	Charset   string  `json:"charset"`
	Steps     []*Step `json:"steps"`
}

type StaticRequest struct {
	TaskID string `json:"task_id"`
}
