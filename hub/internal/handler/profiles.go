package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/api"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
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

	var profileList []protocol.Profile

	switch value := instance.(type) {
	case *rss3uri.PlatformInstance:
		profileList, err = getPlatformInstanceProfileList(value)
		if err != nil {
			_ = c.Error(api.ErrorDatabaseError)

			return
		}
	case *rss3uri.NetworkInstance:
		// TODO
	default:
		_ = c.Error(api.ErrorInvalidParams)

		return
	}

	request := GetProfileListRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(errors.New("invalid params"))

		return
	}

	c.JSON(http.StatusOK, protocol.File{
		// TODO
		DateUpdated: time.Now(),
		Identifier:  fmt.Sprintf("%s/profiles", instance.String()),
		Total:       len(profileList),
		List:        profileList,
	})
}

func getPlatformInstanceProfileList(instance *rss3uri.PlatformInstance) ([]protocol.Profile, error) {
	tx := database.DB.Begin()
	defer tx.Rollback()

	profileModels, err := database.QueryProfiles(tx, instance.GetIdentity(), instance.Platform.ID().Int())
	if err != nil {
		return nil, err
	}

	var profiles []protocol.Profile

	for _, profileModel := range profileModels {
		var attachments []protocol.ProfileAttachment

		for _, attachmentModel := range profileModel.Attachments {
			attachments = append(attachments, protocol.ProfileAttachment{
				Type:     attachmentModel.Type,
				Content:  attachmentModel.Content,
				MimeType: attachmentModel.MimeType,
			})
		}

		var accountConnects []string

		accountModels, err := database.QueryAccounts(
			tx, instance.GetIdentity(), instance.Platform.ID().Int(), constants.ProfileSourceIDCrossbell.Int(),
		)
		if err != nil {
			return nil, err
		}

		for _, accountModel := range accountModels {
			accountConnect := fmt.Sprintf(
				"%s@%s",
				accountModel.ID,
				constants.PlatformID(accountModel.Platform).Symbol().String(),
			)

			accountConnects = append(accountConnects, accountConnect)
		}

		profiles = append(profiles, protocol.Profile{
			DateCreated:       profileModel.CreatedAt,
			DateUpdated:       profileModel.UpdatedAt,
			Name:              database.UnwrapNullString(profileModel.Name),
			Avatars:           profileModel.Avatars,
			Bio:               database.UnwrapNullString(profileModel.Bio),
			Attachments:       attachments,
			ConnectedAccounts: nil,
			Source:            constants.PlatformID(profileModel.Platform).Symbol().String(),
			// TODO
			//Metadata:          protocol.ProfileMetadata{},
		})
	}

	tx.Commit()

	return profiles, nil
}
