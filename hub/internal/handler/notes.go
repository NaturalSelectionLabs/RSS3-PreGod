package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/indexer"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	m "github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/service"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type GetNoteListRequest struct {
	Limit          int      `form:"limit"`
	LastIdentifier string   `form:"last_identifier"`
	Tags           []string `form:"tags"`
	ExcludeTags    []string `form:"exclude_tags"`
	MimeTypes      []string `form:"mime_types"`
	ItemSources    []string `form:"item_sources"`
	Networks       []string `form:"networks"`
	Latest         bool     `form:"latest"`
}

func GetNoteListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	request := GetNoteListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	noteModels, total, err := getNoteListByInstance(c, instance, request)

	if err != nil {
		api.SetError(c, api.ErrorIndexer, err)

		return
	}

	noteList, dateUpdated, errType, err := service.FormatProtocolItemByNote(noteModels)
	if err != nil {
		api.SetError(c, errType, err)

		return
	}

	var lastItem *protocol.Item

	if len(noteList) > 0 {
		lastItem = &noteList[len(noteList)-1]
	}

	identifierNext := ""
	uri := rss3uri.New(instance)

	if total > int64(request.Limit) {
		nextQuery := c.Request.URL.Query()

		if lastItem != nil {
			nextQuery.Set("last_identifier", lastItem.Identifier)
		}

		identifierNext = fmt.Sprintf("%s/notes?%s", uri.String(), nextQuery.Encode())
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated:    dateUpdated,
		Identifier:     fmt.Sprintf("%s/notes?%s", uri.String(), c.Request.URL.Query().Encode()),
		IdentifierNext: identifierNext,
		Total:          total,
		List:           noteList,
	})
}

// nolint:funlen,gocognit // TODO
func getNoteListByInstance(c *gin.Context, instance rss3uri.Instance, request GetNoteListRequest) ([]model.Note, int64, error) {
	if len(request.LastIdentifier) == 0 {
		if err := indexer.GetItems(c.Request.URL.String(), instance, request.Latest); err != nil {
			return nil, 0, err
		}
	}

	// Get instance's notes
	internalDB := database.DB

	if request.LastIdentifier != "" {
		var lastItem model.Note
		if err := database.DB.Where(&model.Note{
			Identifier: strings.ToLower(request.LastIdentifier),
		}).First(&lastItem).Error; err != nil {
			return nil, 0, err
		}

		internalDB = internalDB.
			Where("date_created <= ?", lastItem.DateCreated).
			Where("identifier != ?", lastItem.Identifier).
			Where(
				"(transaction_hash != ?) OR (transaction_hash = ? AND transaction_log_index < ?)",
				lastItem.TransactionHash, lastItem.TransactionHash, lastItem.TransactionLogIndex,
			)
	}

	if request.Tags != nil && len(request.Tags) != 0 {
		internalDB = internalDB.Where("tags && ?", pq.StringArray(request.Tags))
	}

	if request.ExcludeTags != nil && len(request.ExcludeTags) != 0 {
		internalDB = internalDB.Where("tags && ? = FALSE", pq.StringArray(request.ExcludeTags))
	}

	if request.ItemSources != nil && len(request.ItemSources) != 0 {
		internalDB = internalDB.Where("source IN ?", request.ItemSources)
	}

	if request.Networks != nil && len(request.Networks) != 0 {
		internalDB = internalDB.Where("metadata_network IN ?", request.Networks)
	}

	notes := make([]model.Note, 0)
	if err := internalDB.
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Limit(request.Limit).
		Order("date_created DESC").
		Order("transaction_hash DESC").
		Order("transaction_log_index DESC").
		Find(&notes).Error; err != nil {
		return nil, 0, err
	}

	var count int64

	if err := internalDB.
		Model(&model.Note{}).
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Order("date_created DESC").
		Order("transaction_hash DESC").
		Order("transaction_log_index DESC").
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return notes, count, nil
}

// BatchGetNoteListHandlerFunc can batch query notes by request body.
func BatchGetNoteListHandlerFunc(c *gin.Context) {
	req := m.BatchGetNodeListRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	if req.Limit <= 0 {
		req.Limit = middleware.DefaultListLimit
	}

	if req.Limit > middleware.MaxListLimit {
		req.Limit = middleware.MaxListLimit
	}

	resp, errType, err := service.BatchGetNodeList(req)
	if err != nil {
		api.SetError(c, errType, err)

		return
	}

	c.JSON(http.StatusOK, resp)
}
