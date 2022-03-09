package api

import (
	"fmt"
	"net/http"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/status"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/web"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
)

type GetLinkListRequest struct {
	Instance  string `uri:"instance" binding:"required"`
	LinkType  string `uri:"link_type" binding:"required"`
	PageIndex int    `uri:"page_index" binding:"required"`
}

func GetLinkListHandlerFunc(c *gin.Context) {
	request := GetLinkListRequest{}
	if err := c.ShouldBindUri(&request); err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.INVALID_PARAMS, nil)

		return
	}

	instance, err := rss3uri.ParseInstance(request.Instance)
	if err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusBadRequest, status.INVALID_PARAMS, nil)

		return
	}

	// TODO Check if the account exists

	linkListFile := protocol.LinkListFile{
		SignedBase: protocol.SignedBase{
			Base: protocol.Base{
				Version: protocol.Version,
				// TODO Refine rss3uri package
				Identifier: fmt.Sprintf("%s/list/link/following/%d", rss3uri.New(instance).String(), request.PageIndex),
				// TODO IdentifierNext
				// TODO No test data available
				// DateCreated: "",
				// DateUpdated: "",
			},
			Signature: "",
		},
	}

	var links []model.Link
	if err := db.DB.Where(
		"rss3_id = ? and page_index = ?",
		fmt.Sprintf("%s@%s", instance.GetIdentity(), instance.GetSuffix()),
		request.PageIndex,
	).Find(&links).Error; err != nil {
		w := web.Gin{C: c}
		w.JSONResponse(http.StatusInternalServerError, status.ERROR, nil)

		return
	}

	// TODO No test data available
	//for _, _ = range links {
	//	linkListFile.List = append(linkListFile.List, protocol.LinkListFileItem{
	//		Type:             constants.LinkTypeFollowing.String(),
	//		IdentifierTarget: "",
	//	})
	//}

	linkListFile.Total = len(linkListFile.List)

	c.JSON(http.StatusOK, linkListFile)
}
