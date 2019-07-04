package glob

import (
	"errors"
	"fmt"
)

const (
	ScrapeSuccess = "爬取成功"
	ScrapeFailed  = "爬取失敗"
)

func StatusDetail(format string, args ...interface{}) error {
	if len(args) > 0 {
		return errors.New(fmt.Sprintf(format, args...))
	} else {
		return errors.New(format)
	}
}
