package service

import (
	"fmt"
	"time"

	sdk "gitlab.paradise-soft.com.tw/glob/legion-sdk"
)

var (
	ExpiredTime = time.Minute * 5
)

// Notice 透過 kafka 通知使用者取資料 所傳遞的資訊
type Notice sdk.Notice

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RetrySetting sdk.RetrySetting

type LegionRequest sdk.LegionRequest

type LegionResponse sdk.LegionResponse

type LegionResult sdk.LegionResult

type CacheRequest struct {
	Key string `form:"key"`
}

func (r *LegionRequest) CheckKafka() (err error) {
	if r.RespTopic == "" {
		err = fmt.Errorf("RespTopic is empty")
		return
	}

	if r.RawURL == "" {
		err = fmt.Errorf("RawURL is empty")
		return
	}
	return
}
