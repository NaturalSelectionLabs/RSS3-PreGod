package ens

import "time"

type ENSTextRecord struct {
	domain    string
	text      map[string]string
	createdAt time.Time
	txHash    string
}
