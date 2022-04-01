package arweave

import "fmt"

// MirrorContent stores all indexed articles from arweave.
type MirrorContent struct {
	Title          string
	Timestamp      int64
	Content        string
	Author         string
	Link           string
	Digest         string
	OriginalDigest string
	TxHash         string
}

func (a MirrorContent) String() string {
	return fmt.Sprintf(`Title: %s, Timestamp: %d, Author: %s, Link: %s, Digest: %s, OriginalDigest: %s, TxHash: %s`,
		a.Title, a.Timestamp, a.Author, a.Link, a.Digest, a.OriginalDigest, a.TxHash)
}
