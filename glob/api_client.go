package glob

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

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

	// var reader io.Reader
	// switch charset {
	// case "utf-8", "utf8":
	// 	reader = resp.Body
	// case "gbk", "gb18030":
	// 	reader = simplifiedchinese.GB18030.NewDecoder().Reader(resp.Body)
	// default:
	// 	reader = resp.Body
	// }

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
