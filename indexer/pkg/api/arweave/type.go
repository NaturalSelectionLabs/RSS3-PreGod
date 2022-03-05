package arweave

import "fmt"

// MirrorArticle stores all indexed articles from arweave.
type MirrorArticle struct {
	Title          string
	TimeStamp      int64
	Content        string
	Author         string
	Link           string
	Digest         string
	OriginalDigest string
}

func (a MirrorArticle) String() string {
	return fmt.Sprintf(`Title: %s, TimeStamp: %d, Author: %s, Link: %s, Digest: %s, OriginalDigest: %s`,
		a.Title, a.TimeStamp, a.Author, a.Link, a.Digest, a.OriginalDigest)
}
