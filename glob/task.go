package glob

import (
	"time"
)

type Task struct {
	Para      interface{}
	Status    int
	StartTime time.Time
}
