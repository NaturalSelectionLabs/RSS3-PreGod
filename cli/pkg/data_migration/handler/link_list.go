package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/cli/pkg/data_migration/protocol"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
)

//nolint:funlen // no need to split
func LinkList(content bson.D) error {
	// handle link list
	var linkList protocol.RSS3Links031
	// Unmarshal
	doc, err := bson.Marshal(content)
	if err != nil {
		return err
	}

	if err = bson.Unmarshal(doc, &linkList); err != nil {
		return err
	}
	// Split & save into db

	splits := strings.Split(linkList.ID, "-") // {account} - list - links.{id} - {pageNumber}

	account := splits[0]
	listID := strings.Split(splits[2], ".")[1]
	linkTypeId := constants.StringToLinkTypeID(listID)
	pageNumberStr := splits[3]

	instanceID := account + "@" + string(constants.PlatformSymbolEthereum)
	linkListID := instanceID + "/list/links/" + listID + "/" + pageNumberStr

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		return err
	}

	linkListModel := model.LinkList{
		LinkListID: linkListID,
		RSS3ID:     instanceID,
		LinkType:   linkTypeId,
	}

	// Generate link list
	var links []model.Link = make([]model.Link, len(linkList.List))

	for i, target := range linkList.List {
		uniqueName :=
			fmt.Sprint(linkTypeId) +
				"-" + account + "-" + fmt.Sprint(constants.PrefixIDAccount) + "-" + fmt.Sprint(constants.PlatformIDEthereum) +
				"-" + target + "-" + fmt.Sprint(constants.PrefixIDAccount) + "-" + fmt.Sprint(constants.PlatformIDEthereum)

		linkID := uuid.NewV5(uuid.NamespaceOID, uniqueName).String()
		links[i] = model.Link{
			LinkID:           linkID,
			RSS3ID:           account,
			PrefixID:         constants.PrefixIDAccount,
			PlatformID:       constants.PlatformIDEthereum,
			TargetRSS3ID:     target,
			TargetPrefixID:   constants.PrefixIDAccount,
			TargetPlatformID: constants.PlatformIDEthereum,
			PageIndex:        pageNumber,
		}
	}

	// deduplicate links array by linkID
	var dedupedLinks []model.Link

	for _, link := range links {
		var found bool = false

		for _, l := range dedupedLinks {
			if l.LinkID == link.LinkID {
				found = true

				break
			}
		}

		if !found {
			dedupedLinks = append(dedupedLinks, link)
		}
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&linkListModel).Error; err != nil {
			return err
		}

		if len(dedupedLinks) > 0 {
			// insert in batch
			chunkSize := 1000
			for i := 0; i < len(dedupedLinks); i += chunkSize {
				end := i + chunkSize
				if end > len(dedupedLinks) {
					end = len(dedupedLinks)
				}

				if err := tx.Create(dedupedLinks[i:end]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
