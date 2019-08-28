package glob

import (
	"gitlab.paradise-soft.com.tw/glob/dispatcher"
)

func initDispatcher() {
	dispatcher.Init(
		Config.Dispatcher.Brokers,
		dispatcher.InitSetDefaultGroupID(Config.Dispatcher.GroupID),
	)
}
