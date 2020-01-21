module gitlab.paradise-soft.com.tw/dwh/legion

go 1.12

//replace gitlab.paradise-soft.com.tw/glob/legion-sdk => /Users/george_liu/Desktop/golang/common/legion-sdk

//replace gitlab.paradise-soft.com.tw/glob/gorequest => /Users/george_liu/Desktop/golang/common/gorequest

require (
	github.com/DeanThompson/ginpprof v0.0.0-20190408063150-3be636683586
	github.com/PuerkitoBio/goquery v1.5.0
	github.com/chromedp/cdproto v0.0.0-20191114225735-6626966fbae4
	github.com/chromedp/chromedp v0.5.2
	github.com/gin-gonic/gin v1.4.0
	github.com/google/uuid v1.1.1
	github.com/spf13/viper v1.4.0
	gitlab.paradise-soft.com.tw/dwh/proxy v1.2.1
	gitlab.paradise-soft.com.tw/glob/dispatcher v1.11.2
	gitlab.paradise-soft.com.tw/glob/gorequest v0.1.0
	gitlab.paradise-soft.com.tw/glob/helper v0.0.0-20190523032655-2a9a0b97690a
	gitlab.paradise-soft.com.tw/glob/legion-sdk v0.1.10
	gitlab.paradise-soft.com.tw/glob/tracer v1.1.0
	golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
	golang.org/x/text v0.3.0
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
)
