package glob

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	ScrapeSuccess = "爬取成功"
	ScrapeFailed  = "爬取失敗"
)

type ApiClient struct {
	Url        string
	Appkey     string
	Token      string
	TokenLogin string
}

func NewClient(url string) *ApiClient {
	c := &ApiClient{Url: url}
	return c
}

func (c *ApiClient) Get(data map[string]string) ([]byte, error) {
	form := url.Values{}
	addr := c.Url

	if data == nil {
		data = make(map[string]string)
	}

	for k, v := range data {
		form.Add(k, v)
	}

	if len(form) > 0 {
		if strings.IndexAny(addr, "?") > -1 {
			addr += "&" + form.Encode()
		} else {
			addr += "?" + form.Encode()
		}
	}

	return c.doRequest("GET", addr, nil, nil)
}

func (c *ApiClient) doRequest(method, addr string, data map[string]string, header map[string]string) ([]byte, error) {
	dialer := (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	})

	transport := http.Transport{
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := http.Client{
		Transport: &transport,
	}

	form := url.Values{}

	if data != nil {
		for k, v := range data {
			form.Add(k, v)
		}
	}

	req, err := http.NewRequest(method, addr, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	defer func() { req.Close = true }()

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36")
	req.Header.Set("Connection", "close")

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return b, nil
	} else {
		return nil, errors.New(string(b))
	}
}

func StatusDetail(format string, args ...interface{}) error {
	if len(args) > 0 {
		return errors.New(fmt.Sprintf(format, args...))
	} else {
		return errors.New(format)
	}
}
