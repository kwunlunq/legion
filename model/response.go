package model

type Response struct {
	TaskID string `json:"task_id"`
	Body   []byte `json:"body"`
	Error  error  `json:"error"`
}

type StaticResponse struct {
	TaskID  string `json:"task_id"`
	Content string `json:"body"`
}
