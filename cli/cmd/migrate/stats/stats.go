package stats

import (
	"log"
	"sync/atomic"
	"time"
)

var (
	AccountNumber         int64 = 0
	AccountPlatformNumber int64 = 0
	LinkNumber            int64 = 0
	SignatureNumber       int64 = 0
)

func Run() {
	for {
		log.Println(
			"INFO",
			"Account", atomic.LoadInt64(&AccountNumber),
			"AccountPlatform", atomic.LoadInt64(&AccountPlatformNumber),
			"Link", atomic.LoadInt64(&LinkNumber),
			"Signature", atomic.LoadInt64(&SignatureNumber),
		)
		time.Sleep(time.Second)
	}
}
