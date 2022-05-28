package indexer

type Response struct {
	Error struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
}
