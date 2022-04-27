package service

import (
	"encoding/json"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/dao"
	m "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
)

// BatchGetNodeList
// parse the address list into instance list
// query database
// format data
func BatchGetNodeList(req m.BatchGetNodeListRequest) (protocol.File, error, error) {
	req.InstanceList = []rss3uri.Instance{}
	for _, address := range req.AddressList {
		uri, err := rss3uri.Parse(address)
		if err != nil {
			continue
		}
		req.InstanceList = append(req.InstanceList, uri.Instance)
	}
	if len(req.InstanceList) == 0 {
		return protocol.File{}, nil, nil
	}

	noteList, total, err := dao.BatchGetNodeList(req)
	if err != nil {
		return protocol.File{}, api.ErrorDatabase, err
	}

	var dateUpdated *timex.Time
	var itemList = []protocol.Item{}
	for _, note := range noteList {
		attachmentList := []protocol.ItemAttachment{}
		if err = json.Unmarshal(note.Attachments, &attachmentList); err != nil {
			return protocol.File{}, api.ErrorInvalidParams, err
		}

		updated := timex.Time(note.DateCreated)
		if dateUpdated == nil || dateUpdated.Time().Before(note.DateUpdated) {
			dateUpdated = &updated
		}

		metadata := make(map[string]interface{})
		if err = json.Unmarshal(note.Metadata, &metadata); err != nil {
			return protocol.File{}, api.ErrorIndexer, err
		}
		metadata["network"] = note.MetadataNetwork
		metadata["proof"] = note.MetadataProof

		itemList = append(itemList, protocol.Item{
			Identifier:  note.Identifier,
			DateCreated: timex.Time(note.DateCreated),
			DateUpdated: timex.Time(note.DateUpdated),
			RelatedURLs: note.RelatedURLs,
			Links:       fmt.Sprintf("%s/links", note.Identifier),
			BackLinks:   fmt.Sprintf("%s/backlinks", note.Identifier),
			Tags:        note.Tags,
			Authors:     note.Authors,
			Title:       note.Title,
			Summary:     note.Summary,
			Attachments: attachmentList,
			Source:      note.Source,
			Metadata:    metadata,
		})
	}

	file := protocol.File{
		DateUpdated: dateUpdated,
		Total:       total,
		List:        itemList,
	}
	return file, nil, nil
}
