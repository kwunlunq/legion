package service

import (
	"time"

	"gitlab.paradise-soft.com.tw/backend/legion/glob"
	"gitlab.paradise-soft.com.tw/backend/yaitoo/tracer"
)

func CPUMonitor() {
	for {
		cpu, _ := glob.CPUPercent()
		if cpu > float64(glob.Config.CrawlerSetting.CPULimit) {
			tracer.Tracef("CPU", "cpu is up to %g", cpu)
		}
		time.Sleep(30 * time.Second)
	}
}
