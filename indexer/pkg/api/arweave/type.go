package arweave

import "fmt"

// MirrorContent stores all indexed articles from arweave.
type MirrorContent struct {
	Title          string
	TimeStamp      int64
	Content        string
	Author         string
	Link           string
	Digest         string
	OriginalDigest string
	TxHash         string
}

func (a MirrorContent) String() string {
	return fmt.Sprintf(`Title: %s, TimeStamp: %d, Author: %s, Link: %s, Digest: %s, OriginalDigest: %s, TxHash: %s`,
		a.Title, a.TimeStamp, a.Author, a.Link, a.Digest, a.OriginalDigest, a.TxHash)
}
