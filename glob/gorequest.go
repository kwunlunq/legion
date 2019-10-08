package glob

import "gitlab.paradise-soft.com.tw/glob/gorequest"

func NewDefaultGoReq() (goReq *gorequest.SuperAgent) {
	goReq = gorequest.New()
	goReq.SetCurlCommand(Config.Log.Level == "debug")

	// User-Agent from chrome
	goReq.Set(`User-Agent`, `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36`)

	return
}
