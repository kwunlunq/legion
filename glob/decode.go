package glob

import (
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func Decoder(rawContent []byte, encoding string) (result []byte, err error) {
	strings.ToLower(encoding)

	switch encoding {
	case "gbk", "gb18030":
		result, err = simplifiedchinese.GB18030.NewDecoder().Bytes(rawContent)
	case "utf8", "utf-8":
		result = rawContent
	default:
		result = rawContent
	}

	return
}
