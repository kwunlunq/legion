package model

const (
	Click       = "click"
	DoubleClick = "double_cilck"
	SendKeys    = "send_keys"
	WaitReady   = "wait_ready"
	WaitVisible = "wait_visible"
)

type Step struct {
	Action string `json:"action"`
	Target string `json:"target"`
	Keys   string `json:"keys"`
}
