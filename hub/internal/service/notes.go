package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/dao"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/indexer"
	m "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
)

// BatchGetNodeList
// parse the address list into instance list
// query database
// format data
func BatchGetNodeList(ctx context.Context, req m.BatchGetNodeListRequest) (protocol.File, error, error) {
	req.InstanceList = []rss3uri.Instance{}
	for _, address := range req.AddressList {
		uri, err := rss3uri.Parse(address)
		if err != nil {
			continue
		}

		req.InstanceList = append(req.InstanceList, uri.Instance)

		// get item
		if len(req.LastIdentifier) == 0 {
			if err := indexer.GetItems("batch_get_node_list", uri.Instance, req.Latest); err != nil {
				return protocol.File{
					List: make([]protocol.Item, 0),
				}, api.ErrorIndexer, err
			}
		}
	}

	if len(req.InstanceList) == 0 {
		return protocol.File{
			List: make([]protocol.Item, 0),
		}, nil, nil
	}

	noteList, total, err := dao.BatchGetNodeList(req)
	if err != nil {
		return protocol.File{
			List: make([]protocol.Item, 0),
		}, api.ErrorDatabase, err
	}

	itemList, dateUpdated, errType, err := FormatProtocolItemByNote(noteList)
	if err != nil {
		return protocol.File{
			List: make([]protocol.Item, 0),
		}, errType, err
	}

	identifierNext := ""
	if total > int64(req.Limit) {
		identifierNext = itemList[len(itemList)-1].Identifier
	}

	file := protocol.File{
		DateUpdated:    dateUpdated,
		Total:          total,
		List:           itemList,
		IdentifierNext: identifierNext,
	}

	return file, nil, nil
}

// FormatProtocolItemByNote format data
func FormatProtocolItemByNote(noteList []model.Note) ([]protocol.Item, *timex.Time, error, error) {
	var dateUpdated *timex.Time

	var itemList = make([]protocol.Item, 0)

	for _, note := range noteList {
		attachmentList := make([]protocol.ItemAttachment, 0)
		if err := json.Unmarshal(note.Attachments, &attachmentList); err != nil {
			return nil, nil, api.ErrorInvalidParams, err
		}

		updated := timex.Time(note.DateCreated)
		if dateUpdated == nil || dateUpdated.Time().Before(note.DateUpdated) {
			dateUpdated = &updated
		}

		metadata := make(map[string]interface{})
		if err := json.Unmarshal(note.Metadata, &metadata); err != nil {
			return nil, nil, api.ErrorIndexer, err
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

	return itemList, dateUpdated, nil, nil
}
