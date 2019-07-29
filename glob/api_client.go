package glob

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	ServiceSuccess = "服務正常"
	ScrapeSuccess  = "爬取成功"
	ScrapeFailed   = "爬取失敗"
)

func StatusDetail(format string, args ...interface{}) error {
	if len(args) > 0 {
		return errors.New(fmt.Sprintf(format, args...))
	} else {
		return errors.New(format)
	}
}

func GetAndConvertToDocument(targetSite string) (*goquery.Document, error) {
	resp, err := http.Get(targetSite)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	d, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return d, err
}

func GetAndConvertToDocumentByProxy(targetSite string, proxyLocation ...string) (*goquery.Document, error) {
	proxyURL, err := GetProxyErr(proxyLocation...)
	if err != nil {
		return nil, err
	}

	proxy, _ := url.Parse(proxyURL)

	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	resp, err := client.Get(targetSite)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	d, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return d, err
}
