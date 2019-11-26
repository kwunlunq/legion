module gitlab.paradise-soft.com.tw/dwh/legion

go 1.12

//replace gitlab.paradise-soft.com.tw/glob/legion-sdk => /Users/george_liu/Desktop/golang/common/legion-sdk

//replace gitlab.paradise-soft.com.tw/glob/gorequest => /Users/george_liu/Desktop/golang/common/gorequest

require (
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/chromedp/cdproto v0.0.0-20190614062957-d6d2f92b486d
	github.com/chromedp/chromedp v0.3.1-0.20190614072529-35b61282746d
	github.com/gin-gonic/gin v1.4.0
	github.com/google/uuid v1.1.1
	github.com/spf13/viper v1.4.0
	gitlab.paradise-soft.com.tw/dwh/proxy v1.1.0
	gitlab.paradise-soft.com.tw/glob/dispatcher v1.10.12
	gitlab.paradise-soft.com.tw/glob/gorequest v0.1.0
	gitlab.paradise-soft.com.tw/glob/legion-sdk v0.1.2
	gitlab.paradise-soft.com.tw/glob/tracer v1.1.0
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys v0.0.0-20190613124609-5ed2794edfdc // indirect
	golang.org/x/text v0.3.0
)
