package glob

import (
	"gitlab.paradise-soft.com.tw/glob/dispatcher"
)

func initDispatcher() {
	dispatcher.Init(Config.Dispatcher.Brokers, Config.Dispatcher.GroupID)
}
