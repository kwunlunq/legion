package model

type Request struct {
	TaskID        string   `json:"task_id"`
	RespTopic     string   `json:"resp_topic"`
	URL           string   `json:"url"`
	Target        string   `json:"target"`
	ProxyLocation []string `json:"proxy_loc"`
	Charset       string   `json:"charset"`
	Steps         []*Step  `json:"steps"`
	Body          []byte   `json:"body"`
	Error         error    `json:"error"`
}

type CacheRequest struct {
	TaskID  string `json:"task_id" form:"task_id"`
	Content string `json:"content"`
}
