package service

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/glob/gorequest"
	sdk "gitlab.paradise-soft.com.tw/glob/legion-sdk"
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

func (r *LegionRequest) GetStaticResult() (legionResult *LegionResult) {
	var resp *http.Response
	var body []byte
	var err error
	resp, body, err = r.doStatic()

	legionResult = &LegionResult{}
	legionResult.Request = (*sdk.LegionRequest)(r)
	if err != nil {
		legionResult.ErrorMessage = err.Error()
		return
	}

	legionResp := &LegionResponse{}
	legionResp.FinishedAt = time.Now()
	legionResp.StatusCode = resp.StatusCode
	legionResp.Header = make(map[string]string, len(resp.Header))
	for key, val := range resp.Header {
		legionResp.Header[key] = strings.Join(val, ",")
	}

	legionResp.Body = body
	legionResult.Response = (*sdk.LegionResponse)(legionResp)
	return
}

func (r *LegionRequest) toGoRequest() (goReq *gorequest.SuperAgent, err error) {
	goReq = glob.NewDefaultGoReq()
	goReq.Url = r.RawURL
	goReq.QueryData = r.QueryData

	goReq.Method = "GET"
	if r.Method != "" {
		goReq.Method = r.Method
	}

	if r.TargetType != "" {
		goReq.Type(r.TargetType)
	}

	// goReq.Send(r.Body)
	goReq.SendString(string(r.Body))

	for k, v := range r.Header {
		goReq.Set(k, v)
	}

	goReq.AddCookies(r.Cookies)

	if r.BasicAuth != nil {
		goReq.SetBasicAuth(r.BasicAuth.Username, r.BasicAuth.Password)
	}

	if r.Timeout != 0 {
		goReq.Timeout(r.Timeout)
	} else {
		goReq.Timeout(glob.Config.GoRequest.Timeout)
	}

	if len(r.ProxyLocations) > 0 {
		var proxies []string
		proxies, err = glob.GetProxies(len(r.ProxyLocations), r.ProxyLocations)
		if err != nil {
			return nil, err
		}
		goReq.Proxy(proxies...)
	}

	if len(r.Proxies) > 0 {
		goReq.Proxy(r.Proxies...)
	}

	if r.InsecureSkipVerify {
		goReq.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	if r.Retryable != nil {
		goReq.Retry(
			r.Retryable.RetryCount,
			r.Retryable.RetryTime,
			r.Retryable.RetryableStatus...,
		)
	}

	return goReq, nil
}

func (r *LegionRequest) doStatic() (resp *http.Response, body []byte, err error) {
	defer func() {
		if err != nil {
			resp = nil
			body = nil
		}
	}()

	goReq, err := r.toGoRequest()
	if err != nil {
		return
	}

	var errs glob.Errors
	resp, body, errs = goReq.EndBytes()
	if !errs.IsNil() {
		err = errs
		return
	}

	if r.Target != "" {
		var goDoc *goquery.Document
		goDoc, err = goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return
		}
		body = []byte(goDoc.Find(r.Target).Text())
	}

	if r.Charset != "" {
		body, err = glob.Decoder(body, r.Charset)
		if err != nil {
			return
		}
	}

	return
}
