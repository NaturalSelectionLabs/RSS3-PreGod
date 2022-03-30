package handler

import (
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetProfileListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	profileList := []protocol.Profile{
		{
			DateCreated:       time.Time{},
			DateUpdated:       time.Time{},
			Name:              "",
			Avatars:           nil,
			Bio:               "",
			Attachments:       nil,
			ConnectedAccounts: nil,
			Source:            "",
			Metadata:          protocol.ProfileMetadata{},
		},
	}

	c.JSON(http.StatusOK, protocol.File{
		DateUpdated: time.Now(),
		Identifier:  fmt.Sprintf("%s/profiles", instance.String()),
		Total:       len(profileList),
		List:        profileList,
	})
}
