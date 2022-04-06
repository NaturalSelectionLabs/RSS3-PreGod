package arweave

import "fmt"

type ArLatestBlockResult struct {
	Network          string `json:"network"`
	Version          int64  `json:"version"`
	Release          int64  `json:"release"`
	Height           int64  `json:"height"`
	Current          string `json:"current"`
	Blocks           int64  `json:"blocks"`
	Peers            int64  `json:"peers"`
	QueueLength      int64  `json:"queue_length"`
	NodeStateLatency int64  `json:"node_state_latency"`
}

type GraphqlResultEdges struct {
	Node struct {
		Id   string `json:"id"`
		Tags []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"tags"`
	} `json:"node"`
}

type GraphqlResult struct {
	Data struct {
		Transactions struct {
			Edges []GraphqlResultEdges `json:"edges"`
		} `json:"transactions"`
	} `json:"data"`
}

type OriginalMirrorContent struct {
	Content struct {
		Body      string `json:"body"`
		Timestamp int64  `json:"timestamp"`
		Title     string `json:"title"`
	} `json:"content"`
	Digest     string `json:"digest"`
	Authorship struct {
		Contributor         string `json:"contributor"`
		SigningKey          string `json:"signingKey"` // nolint:tagliatelle // returned by arewave api
		Signature           string `json:"signature"`
		SigningKeySignature string `json:"signingKeySignature"` // nolint:tagliatelle // returned by arewave api
		SigningKeyMessage   string `json:"signingKeyMessage"`   // nolint:tagliatelle // returned by arewave api
		Algorithm           struct {
			Name string `json:"name"`
			Hash string `json:"hash"`
		} `json:"algorithm"`
	} `json:"authorship"`
	Nft struct {
	} `json:"nft"`
	Version        string `json:"version"`
	OriginalDigest string `json:"originalDigest"` // nolint:tagliatelle // returned by arewave api
}

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
