package indexer

import "time"

const (
	Endpoint     = "http://localhost:8081"
	EndpointItem = "http://localhost:8081/item"
)

type Response struct {
	Error struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
	Data struct {
		Note []struct {
			ItemId struct {
				NetworkId int    `json:"network_id"`
				Proof     string `json:"proof"`
			} `json:"item_id"`
			Metadata struct {
			} `json:"metadata"`
			Tags        []string `json:"tags"`
			Authors     []string `json:"authors"`
			Title       string   `json:"title"`
			Summary     string   `json:"summary"`
			Attachments []struct {
				Content     string    `json:"content"`
				Address     string    `json:"address"`
				MimeType    string    `json:"mime_type"`
				SizeInBytes int64     `json:"size_in_bytes"`
				SyncAt      time.Time `json:"sync_at"`
			} `json:"attachments"`
			DateCreated string `json:"date_created"`
		} `json:"note"`
		Asset []struct {
			ItemId struct {
				NetworkId int    `json:"network_id"`
				Proof     string `json:"proof"`
			} `json:"item_id"`
			Metadata struct {
			} `json:"metadata"`
			Tags        []string `json:"tags"`
			Authors     []string `json:"authors"`
			Title       string   `json:"title"`
			Summary     string   `json:"summary"`
			Attachments []struct {
				Content     string    `json:"content"`
				Address     string    `json:"address"`
				MimeType    string    `json:"mime_type"`
				SizeInBytes int64     `json:"size_in_bytes"`
				SyncAt      time.Time `json:"sync_at"`
			} `json:"attachments"`
			DateCreated string `json:"date_created"`
		} `json:"asset"`
	} `json:"data"`
}
