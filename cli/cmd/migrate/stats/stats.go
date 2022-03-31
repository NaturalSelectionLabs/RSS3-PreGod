package stats

import (
	"log"
	"sync/atomic"
	"time"
)

var (
	Profile int64 = 0
	Account int64 = 0
	Link    int64 = 0
)

func Run() {
	for {
		log.Println(
			"INFO",
			"Profile", atomic.LoadInt64(&Profile),
			"Account", atomic.LoadInt64(&Account),
			"Link", atomic.LoadInt64(&Link),
		)
		time.Sleep(time.Second)
	}
}
