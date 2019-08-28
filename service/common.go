package service

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	ExpiredTime = time.Minute * 1
)

// Notice 透過 kafka 通知使用者取資料 所傳遞的資訊
type Notice struct {
	InternalURL string    `json:"internalURL"`
	ExternalURL string    `json:"externalURL"`
	CreatedAt   time.Time `json:"createdAt"`
}

type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Retryable struct {
	RetryerCount    int           `json:"retryerCount"`
	RetryerTime     time.Duration `json:"retryerTime"`
	RetryableStatus []int         `json:"retryableStatus"`
}

type LegionRequest struct {
	SentAt    time.Time   `json:"sentAt"` // for kafka timeout task
	RespTopic string      `json:"respTopic"`
	Data      interface{} `json:"data"` // Data 使用者傳遞資料 會回傳給使用者

	Steps   []*Step `json:"steps"`
	Target  string  `json:"target"`
	Charset string  `json:"charset"`

	RawURL             string            `json:"rawURL" binding:"required"`
	QueryData          url.Values        `json:"queryData"` // 可能不用
	Method             string            `json:"method"`
	TargetType         string            `json:"targetType"`
	Body               []byte            `json:"body"`
	Header             map[string]string `json:"header"`
	Cookies            []*http.Cookie    `json:"cookies"`
	BasicAuth          *BasicAuth        `json:"basicAuth"`
	Timeout            time.Duration     `json:"timeout"`
	InsecureSkipVerify bool              `json:"insecureSkipVerify"`
	Proxies            []string          `json:"proxies"`        // 每次嘗試將採用的 proxy, "" 代表 local
	ProxyLocations     []string          `json:"proxyLocations"` // 每次嘗試將採用的 proxy 以國家代碼表示 如 cn us, "" 代表 local
	Retryable          *Retryable        `json:"retryable"`
}

type LegionResponse struct {
	StatusCode int               `json:"statusCode"`
	Header     map[string]string `json:"header"`
	Body       []byte            `json:"responseBody"`
}

type LegionResult struct {
	ErrorMessage string          `json:"errorMessage"`
	Request      *LegionRequest  `json:"request"`
	Response     *LegionResponse `json:"response"`
}

type CacheRequest struct {
	Key string `form:"key"`
}

func (this *LegionRequest) CheckKafka() (err error) {
	if this.RespTopic == "" {
		err = fmt.Errorf("RespTopic is empty")
		return
	}

	if this.RawURL == "" {
		err = fmt.Errorf("RawURL is empty")
		return
	}
	return
}
