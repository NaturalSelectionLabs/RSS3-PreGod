package handler

import (
	"encoding/json"
	"errors"
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
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/timex"
	"github.com/gin-gonic/gin"
)

type GetProfileListRequest struct {
	ProfileSources []string `form:"profile_sources"`
}

func GetProfileListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetInstance(c)
	if err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	request := GetProfileListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		api.SetError(c, api.ErrorInvalidParams, err)

		return
	}

	var profileList []protocol.Profile

	var total int64

	switch value := instance.(type) {
	case *rss3uri.PlatformInstance:
		profileList, total, err = getPlatformInstanceProfileList(value, request)
		if err != nil {
			api.SetError(c, api.ErrorIndexer, err)

			return
		}
	case *rss3uri.NetworkInstance:
		switch value.Prefix {
		case constants.PrefixNameAsset:
			profileList, total, err = getAssetProfile(value, request)
			if err != nil {
				api.SetError(c, api.ErrorIndexer, err)

				return
			}
		default:
			api.SetError(c, api.ErrorInvalidParams, errors.New("unsupported prefix name"))

			return
		}
	default:
		api.SetError(c, api.ErrorInvalidParams, errors.New("unsupported instance type"))

		return
	}

	var dateUpdated *timex.Time

	for _, profile := range profileList {
		internalTime := profile.DateUpdated
		if dateUpdated == nil {
			dateUpdated = &internalTime
		} else if dateUpdated.Time().Before(profile.DateUpdated.Time()) {
			dateUpdated = &internalTime
		}
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated: dateUpdated,
		Identifier:  fmt.Sprintf("%s/profiles?%s", rss3uri.New(instance), c.Request.URL.Query().Encode()),
		Total:       total,
		List:        profileList,
	})
}

// nolint:funlen // TODO
func getPlatformInstanceProfileList(instance *rss3uri.PlatformInstance, request GetProfileListRequest) ([]protocol.Profile, int64, error) {
	var profileModels []model.Profile

	internalDB := database.DB

	internalDB = internalDB.Where(&model.Profile{
		ID:       instance.GetIdentity(),
		Platform: instance.Platform.ID().Int(),
	})

	if request.ProfileSources != nil && len(request.ProfileSources) > 0 {
		profileSources := make([]int, len(request.ProfileSources))

		for i, source := range request.ProfileSources {
			profileSources[i] = constants.ProfileSourceName(source).ID().Int()
		}

		internalDB = internalDB.Where("source IN ?", profileSources)
	}

	if err := internalDB.Find(&profileModels).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := internalDB.Model(&model.Profile{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	profiles := make([]protocol.Profile, 0)

	for _, profileModel := range profileModels {
		attachments := make([]protocol.ProfileAttachment, 0)

		if profileModel.Attachments != nil && len(profileModel.Attachments) > 0 {
			if err := json.Unmarshal(profileModel.Attachments, &attachments); err != nil {
				return nil, 0, err
			}
		}

		connectedAccounts := make([]string, 0)

		accountModels, err := database.QueryAccounts(
			database.DB,
			instance.GetIdentity(),
			instance.Platform.ID().Int(),
			// TODO
			constants.ProfileSourceIDCrossbell.Int(),
		)
		if err != nil {
			return nil, 0, err
		}

		for _, accountModel := range accountModels {
			connectedAccounts = append(
				connectedAccounts,
				rss3uri.New(rss3uri.NewAccountInstance(accountModel.Identity, constants.PlatformID(accountModel.Platform).Symbol())).String(),
			)
		}

		profiles = append(profiles, protocol.Profile{
			DateCreated:       timex.Time(profileModel.CreatedAt),
			DateUpdated:       timex.Time(profileModel.UpdatedAt),
			Name:              database.UnwrapNullString(profileModel.Name),
			Avatars:           profileModel.Avatars,
			Bio:               database.UnwrapNullString(profileModel.Bio),
			Attachments:       attachments,
			ConnectedAccounts: connectedAccounts,
			Source:            constants.ProfileSourceID(profileModel.Source).Name().String(),
			Metadata: protocol.ProfileMetadata{
				// TODO Now only Crossbell is supported
				Network: strings.ToLower(constants.NetworkSymbolCrossbell.String()),
				Proof:   instance.Identity,
			},
		})

		go func() {
			if err := indexer.GetItems(instance, accountModels); err != nil {
				logger.Error(err)
			}
		}()
	}

	return profiles, count, nil
}

func getAssetProfile(instance *rss3uri.NetworkInstance, request GetProfileListRequest) ([]protocol.Profile, int64, error) {
	internalDB := database.DB

	if request.ProfileSources != nil && len(request.ProfileSources) != 0 {
		profileSources := make([]int, 0)

		for _, source := range request.ProfileSources {
			profileSources = append(profileSources, constants.ProfileSourceName(source).ID().Int())
		}

		internalDB = internalDB.Where("source IN ?", profileSources)
	}

	asset := model.Asset{}
	if err := internalDB.Where(&model.Asset{
		Identifier: strings.ToLower(instance.UriString()),
	}).First(&asset).Error; err != nil {
		return nil, 0, err
	}

	attachments := make([]protocol.ProfileAttachment, 0)
	if err := json.Unmarshal(asset.Attachments, &attachments); err != nil {
		return nil, 0, err
	}

	profiles := []protocol.Profile{
		{
			DateCreated: timex.Time(asset.DateCreated),
			DateUpdated: timex.Time(asset.DateUpdated),
			Name:        asset.Title,
			Bio:         asset.Summary,
			Attachments: attachments,
			Source:      asset.Source,
			Metadata: protocol.ProfileMetadata{
				Network: strings.ToLower(asset.MetadataNetwork),
				Proof:   asset.MetadataProof,
			},
		},
	}

	return profiles, int64(len(profiles)), nil
}
