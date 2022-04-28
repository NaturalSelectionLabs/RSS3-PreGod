package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/indexer"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type GetAssetListRequest struct {
	Limit          int      `form:"limit"`
	LastIdentifier string   `form:"last_identifier"`
	Tags           []string `form:"tags"`
	ExcludeTags    []string `form:"exclude_tags"`
	MimeTypes      []string `form:"mime_types"`
	ItemSources    []string `form:"item_sources"`
	LinkSources    []string `form:"link_sources"`
	LinkType       string   `form:"link_type"`
	Networks       []string `form:"networks"`
	ProfileSources []string `form:"profile_sources"`
}

// nolint:funlen // TODO
func GetAssetListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	request := GetAssetListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	var assetModels []model.Asset

	var total int64

	if len(request.LinkSources) != 0 || request.LinkType != "" {
		assetModels, total, err = getAssetListsByLink(c, instance, request)
	} else {
		assetModels, total, err = getAssetListByInstance(c, instance, request)
	}

	if err != nil {
		api.SetError(c, api.ErrorDatabase, err)

		return
	}

	uri := rss3uri.New(instance)

	var dateUpdated *timex.Time

	assetList := make([]protocol.Item, len(assetModels))

	for i, assetModel := range assetModels {
		attachmentList := make([]protocol.ItemAttachment, 0)
		if err = json.Unmarshal(assetModel.Attachments, &attachmentList); err != nil {
			api.SetError(c, api.ErrorIndexer, err)

			return
		}

		internalTime := timex.Time(assetModel.DateUpdated)
		if dateUpdated == nil {
			dateUpdated = &internalTime
		} else if dateUpdated.Time().Before(assetModel.DateUpdated) {
			dateUpdated = &internalTime
		}

		// Build metadata
		metadata := make(map[string]interface{})

		if err := json.Unmarshal(assetModel.Metadata, &metadata); err != nil {
			api.SetError(c, api.ErrorIndexer, err)

			return
		}

		metadata["network"] = assetModel.MetadataNetwork
		metadata["proof"] = assetModel.MetadataProof

		assetList[i] = protocol.Item{
			Identifier:  assetModel.Identifier,
			DateCreated: timex.Time(assetModel.DateCreated),
			DateUpdated: timex.Time(assetModel.DateUpdated),
			RelatedURLs: assetModel.RelatedURLs,
			Links:       fmt.Sprintf("%s/links", assetModel.Identifier),
			BackLinks:   fmt.Sprintf("%s/backlinks", assetModel.Identifier),
			Tags:        assetModel.Tags,
			Authors:     assetModel.Authors,
			Title:       assetModel.Title,
			Summary:     assetModel.Summary,
			Attachments: attachmentList,
			Source:      assetModel.Source,
			Metadata:    metadata,
		}
	}

	var lastItem *protocol.Item

	for _, item := range assetList {
		internalItem := item

		if lastItem == nil || lastItem.DateCreated.Time().After(internalItem.DateCreated.Time()) {
			lastItem = &internalItem
		}
	}

	identifierNext := ""

	if len(assetList) == database.MaxLimit {
		nextQuery := c.Request.URL.Query()
		if lastItem != nil {
			nextQuery.Set("last_identifier", lastItem.Identifier)
		}

		identifierNext = fmt.Sprintf("%s/assets?%s", uri.String(), nextQuery.Encode())
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated:    dateUpdated,
		Identifier:     fmt.Sprintf("%s/assets?%s", uri.String(), c.Request.URL.Query().Encode()),
		IdentifierNext: identifierNext,
		Total:          total,
		List:           assetList,
	})
}

// nolint:funlen,gocognit // TODO
func getAssetListByInstance(c *gin.Context, instance rss3uri.Instance, request GetAssetListRequest) ([]model.Asset, int64, error) {
	// Get instance's profiles
	var profiles []model.Profile

	internalDB := database.DB

	if request.ProfileSources != nil && len(request.ProfileSources) > 0 {
		var profileSources []int
		for _, source := range request.ProfileSources {
			profileSources = append(profileSources, constants.ProfileSourceName(source).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", profileSources)
	}

	if err := internalDB.Where(&model.Profile{
		ID:       strings.ToLower(instance.GetIdentity()),
		Platform: constants.PlatformSymbol(instance.GetSuffix()).ID().Int(),
	}).Find(&profiles).Error; err != nil {
		return nil, 0, err
	}

	profileIDs := make([]string, len(profiles))
	for _, profile := range profiles {
		profileIDs = append(profileIDs, profile.ID)
	}

	// Get instance's all accounts
	accounts := make([]model.Account, 0)
	if err := internalDB.
		Where("profile_id IN ?", profileIDs).
		Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	if err := indexer.GetItems(c.Request.URL, instance, accounts); err != nil {
		return nil, 0, err
	}

	// Get instance's notes
	internalDB = database.DB

	if request.LastIdentifier != "" {
		var lastItem model.Asset
		if err := database.DB.Where(&model.Asset{
			Identifier: strings.ToLower(request.LastIdentifier),
		}).First(&lastItem).Error; err != nil {
			return nil, 0, err
		}

		internalDB = internalDB.
			Where("date_created <= ?", lastItem.DateCreated).
			Where("identifier != ?", lastItem.Identifier).
			Where("token_id < ?", lastItem.TokenID)
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

	if request.ProfileSources != nil && len(request.ProfileSources) != 0 {
		authors := []string{
			rss3uri.New(instance).String(),
		}

		for _, account := range accounts {
			accountInstance := rss3uri.NewAccountInstance(account.Identity, constants.PlatformID(account.Platform).Symbol())
			authors = append(authors, rss3uri.New(accountInstance).String())
		}

		internalDB = internalDB.Where("authors && ?", pq.StringArray(authors))
	}

	if request.Networks != nil && len(request.Networks) != 0 {
		internalDB = internalDB.Where("metadata_network IN ?", request.Networks)
	}

	assets := make([]model.Asset, 0)
	if err := internalDB.
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Limit(request.Limit).
		Order("date_created DESC").
		Order("identifier DESC").
		Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	var count int64

	if err := internalDB.
		Model(&model.Asset{}).
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Order("date_created DESC").
		Order("identifier DESC").
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return assets, count, nil
}

// nolint:funlen,gocognit // TODO
func getAssetListsByLink(c *gin.Context, instance rss3uri.Instance, request GetAssetListRequest) ([]model.Asset, int64, error) {
	links := make([]model.Link, 0)

	internalDB := database.DB

	if request.ProfileSources != nil && len(request.ProfileSources) > 0 {
		var profileSources []int
		for _, source := range request.ProfileSources {
			profileSources = append(profileSources, constants.ProfileSourceName(source).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", profileSources)
	}

	if request.LinkType != "" {
		internalDB = internalDB.Where("type = ?", constants.LinkTypeName(request.LinkType).ID().Int())
	}

	if request.LinkSources != nil && len(request.LinkSources) != 0 {
		var sources []int
		for _, linkSource := range request.LinkSources {
			sources = append(sources, constants.LinkSourceName(linkSource).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", sources)
	}

	if err := internalDB.
		Where(&model.Link{
			From: instance.GetIdentity(),
		}).
		Find(&links).Error; err != nil {
		return nil, 0, err
	}

	targets := make([]string, 0)

	for _, link := range links {
		targets = append(targets, link.To)
	}

	accounts := make([]model.Account, 0)
	if err := database.DB.
		Where("profile_id IN ?", targets).
		// TODO profile_platform
		Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	if err := indexer.GetItems(c.Request.URL, instance, accounts); err != nil {
		return nil, 0, err
	}

	owners := make([]string, len(links))

	for _, link := range links {
		ownerInstance, err := rss3uri.NewInstance(
			constants.InstanceTypeID(link.ToInstanceType).String(),
			link.To,
			constants.PlatformID(link.ToPlatformID).Symbol().String(),
		)
		if err != nil {
			return nil, 0, err
		}

		owners = append(owners, strings.ToLower(rss3uri.New(ownerInstance).String()))
	}

	internalDB = database.DB

	if request.LastIdentifier != "" {
		var lastItem model.Asset
		if err := database.DB.Where(&model.Asset{
			Identifier: strings.ToLower(request.LastIdentifier),
		}).First(&lastItem).Error; err != nil {
			return nil, 0, err
		}

		internalDB = internalDB.
			Where("date_created <= ?", lastItem.DateCreated).
			Where("identifier != ?", lastItem.Identifier).
			Where("token_id < ?", lastItem.TokenID)
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

	if request.ProfileSources != nil && len(request.ProfileSources) != 0 {
		authors := []string{
			rss3uri.New(instance).String(),
		}

		for _, account := range accounts {
			accountInstance := rss3uri.NewAccountInstance(account.Identity, constants.PlatformID(account.Platform).Symbol())
			authors = append(authors, rss3uri.New(accountInstance).String())
		}

		internalDB = internalDB.Where("authors && ?", pq.StringArray(authors))
	}

	if request.Networks != nil && len(request.Networks) != 0 {
		internalDB = internalDB.Where("metadata_network IN ?", request.Networks)
	}

	assets := make([]model.Asset, 0)
	if err := internalDB.
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Limit(request.Limit).
		Order("date_created DESC").
		Order("identifier DESC").
		Find(&assets).Error; err != nil {
		return nil, 0, err
	}

	var count int64

	if err := internalDB.
		Model(&model.Asset{}).
		Where("owner = ?", strings.ToLower(rss3uri.New(instance).String())).
		Order("date_created DESC").
		Order("identifier DESC").
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return assets, count, nil
}
