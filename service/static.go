package service

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/glob/gorequest"
)

const (
	// TargetType
	TypeJSON      = "json"      // "application/json"
	TypeXML       = "xml"       // "application/xml"
	TypeForm      = "form"      // "application/x-www-form-urlencoded"
	TypeHTML      = "html"      // "text/html"
	TypeText      = "text"      // "text/plain"
	TypeMultipart = "multipart" // "multipart/form-data"
)

func (this *LegionRequest) GetStaticResult() (legionResult *LegionResult) {
	var resp *http.Response
	var body []byte
	var err error
	resp, body, err = this.doStatic()

	legionResult = &LegionResult{}
	legionResult.Request = this
	if err != nil {
		legionResult.ErrorMessage = err.Error()
		return
	}

	legionResp := &LegionResponse{}
	legionResp.StatusCode = resp.StatusCode
	legionResp.Header = make(map[string]string, len(resp.Header))
	for key, val := range resp.Header {
		legionResp.Header[key] = strings.Join(val, ",")
	}

	legionResp.Body = body
	legionResult.Response = legionResp
	return
}

func (this *LegionRequest) toGoRequest() (goReq *gorequest.SuperAgent, err error) {
	goReq = glob.NewDefaultGoReq()
	goReq.Url = this.RawURL
	goReq.QueryData = this.QueryData

	goReq.Method = "GET"
	if this.Method != "" {
		goReq.Method = this.Method
	}

	if this.TargetType != "" {
		goReq.Type(this.TargetType)
	}

	// goReq.Send(this.Body)
	goReq.SendString(string(this.Body))

	for k, v := range this.Header {
		goReq.Set(k, v)
	}

	goReq.AddCookies(this.Cookies)

	if this.BasicAuth != nil {
		goReq.SetBasicAuth(this.BasicAuth.Username, this.BasicAuth.Password)
	}

	if this.Timeout != 0 {
		goReq.Timeout(this.Timeout)
	} else {
		goReq.Timeout(glob.Config.GoRequest.Timeout)
	}

	if len(this.ProxyLocations) > 0 {
		var proxies []string
		proxies, err = glob.GetProxiecErr(len(this.ProxyLocations), this.ProxyLocations)
		if err != nil {
			return nil, err
		}
		this.Proxies = append(this.Proxies, proxies...)
	}

	if len(this.Proxies) > 0 {
		goReq.Proxy(this.Proxies...)
	}

	if this.InsecureSkipVerify {
		goReq.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	if this.Retryable != nil {
		goReq.Retry(
			this.Retryable.RetryerCount,
			this.Retryable.RetryerTime,
			this.Retryable.RetryableStatus...,
		)
	}

	return goReq, nil
}

func (this *LegionRequest) doStatic() (resp *http.Response, body []byte, err error) {
	defer func() {
		if err != nil {
			resp = nil
			body = nil
		}
	}()

	goReq, err := this.toGoRequest()
	if err != nil {
		return
	}

	var errs glob.Errors
	resp, body, errs = goReq.EndBytes()
	if !errs.IsNil() {
		err = errs
		return
	}

	if this.Target != "" {
		var goDoc *goquery.Document
		goDoc, err = goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return
		}
		body = []byte(goDoc.Find(this.Target).Text())
	}

	if this.Charset != "" {
		body, err = glob.Decoder(body, this.Charset)
		if err != nil {
			return
		}
	}

	return
}
