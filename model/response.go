package model

type Response struct {
	TaskID string `json:"task_id"`
	Body   []byte `json:"body"`
	Error  string `json:"error"`
}