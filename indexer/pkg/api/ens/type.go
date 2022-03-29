package ens

import "time"

type ENSTextRecord struct {
	Domain      string
	Description string
	Text        map[string]string
	CreatedAt   time.Time
	TxHash      string
}

// returns a list of recommended keys for a given ENS domain, as per https://app.ens.domains/
// this is a combination of Global Keys and Service Keys, see: https://eips.ethereum.org/EIPS/eip-634
func getTextRecordKeyList() []string {
	// nolint:lll // this is a list of keys
	return []string{"email", "url", "avatar", "description", "notice", "keywords", "com.discord", "com.github", "com.reddit", "com.twitter", "org.telegram", "eth.ens.delegate"}
}
