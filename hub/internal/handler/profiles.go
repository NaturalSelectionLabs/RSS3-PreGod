package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type GetProfileListRequest struct {
	ProfileSources []string `form:"profile_sources"`
}

func GetProfileListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetInstance(c)
	if err != nil {
		_ = c.Error(api.ErrorInvalidParams)

		return
	}

	request := GetProfileListRequest{}
	if err = c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(api.ErrorInvalidParams)

		return
	}

	profileList := make([]protocol.Profile, 0)

	switch value := instance.(type) {
	case *rss3uri.PlatformInstance:
		profileList, err = getPlatformInstanceProfileList(value, request)
		if err != nil {
			if !errors.Is(err, api.ErrorNotFound) {
				err = api.ErrorDatabase
			}

			_ = c.Error(err)

			return
		}
	case *rss3uri.NetworkInstance:
		// TODO
	default:
		_ = c.Error(api.ErrorInvalidParams)

		return
	}

	var dateUpdated *time.Time
	for _, profile := range profileList {
		if dateUpdated == nil {
			// nolint:exportloopref // TODO
			dateUpdated = &profile.DateUpdated
		} else if dateUpdated.Before(profile.DateCreated) {
			// nolint:exportloopref // TODO
			dateUpdated = &profile.DateUpdated
		}
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated: dateUpdated,
		Identifier:  fmt.Sprintf("%s/profiles", rss3uri.New(instance)),
		Total:       len(profileList),
		List:        profileList,
	})
}

func getPlatformInstanceProfileList(instance *rss3uri.PlatformInstance, request GetProfileListRequest) ([]protocol.Profile, error) {
	profileSources := make([]int, 0)
	for _, profileSource := range request.ProfileSources {
		profileSources = append(profileSources, constants.ProfileSourceName(profileSource).ID().Int())
	}

	profileModels, err := database.QueryProfiles(database.DB, instance.GetIdentity(), instance.Platform.ID().Int(), profileSources)
	if err != nil {
		return nil, err
	}

	profiles := make([]protocol.Profile, 0)

	for _, profileModel := range profileModels {
		attachments := make([]protocol.ProfileAttachment, 0)

		for _, attachmentModel := range profileModel.Attachments {
			attachments = append(attachments, protocol.ProfileAttachment{
				Type:     attachmentModel.Type,
				Content:  attachmentModel.Content,
				MimeType: attachmentModel.MimeType,
			})
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
			return nil, err
		}

		for _, accountModel := range accountModels {
			connectedAccounts = append(
				connectedAccounts,
				rss3uri.New(rss3uri.NewAccountInstance(accountModel.ID, constants.PlatformID(accountModel.Platform).Symbol())).String(),
			)
		}

		profiles = append(profiles, protocol.Profile{
			DateCreated:       profileModel.CreatedAt,
			DateUpdated:       profileModel.UpdatedAt,
			Name:              database.UnwrapNullString(profileModel.Name),
			Avatars:           profileModel.Avatars,
			Bio:               database.UnwrapNullString(profileModel.Bio),
			Attachments:       attachments,
			ConnectedAccounts: connectedAccounts,
			Source:            constants.ProfileSourceID(profileModel.Source).Name().String(),
			Metadata: protocol.ProfileMetadata{
				// TODO Now only Crossbell is supported,
				Network: cases.Title(language.English, cases.NoLower).String(constants.NetworkSymbolCrossbell.String()),
				// Network: constants.NetworkID(profileModel.Platform).Symbol().String(),
				Proof: instance.Identity,
			},
		})
	}

	return profiles, nil
}
