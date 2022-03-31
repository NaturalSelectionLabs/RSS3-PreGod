package handler

import (
	"fmt"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/middleware"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/internal/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetBackLinkListHandlerFunc(c *gin.Context) {
	instance, err := middleware.GetPlatformInstance(c)
	if err != nil {
		return
	}

	linkItem := []protocol.LinkItem{
		{
			DateCreated: time.Now(),
			From:        "0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944",
			To:          "0x0fefeD77Bb715E96f1c35c1a4E0D349563d6f6c0",
			Source:      "Corssbell",
			Metadata: protocol.LinkItemMetadata{
				Network: "Crossbell",
				Proof:   "todo",
			},
		},
		{
			DateCreated: time.Now(),
			From:        "0xC8b960D09C0078c18Dcbe7eB9AB9d816BcCa8944",
			To:          "0x0fefeD77Bb715E96f1c35c1a4E0D349563d6f6c0",
			Source:      "Lens",
			Metadata: protocol.LinkItemMetadata{
				Network: "Polygon",
				Proof:   "todo",
			},
		},
	}

	c.JSON(http.StatusOK, protocol.File{
		Identifier:  fmt.Sprintf("%s/backlinks", rss3uri.New(instance).String()),
		DateUpdated: time.Now(),
		Total:       len(linkItem),
		List:        linkItem,
	})
}
