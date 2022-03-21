package database

import (
	"context"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/hub/database/model"
	"gorm.io/gorm"
)

var (
	_ Database = &database{}
)

type database struct {
	db *gorm.DB
}

func (d *database) DB(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx)
}

func (d *database) Tx(ctx context.Context) *gorm.DB {
	return d.db.WithContext(ctx).Begin()
}

func (d *database) QueryAccount(db *gorm.DB, id string, platformID int) (*model.Account, error) {
	account := model.Account{}
	if err := db.Where("id = ? and platform = ?", id, platformID).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (d *database) QueryAccountPlatforms(db *gorm.DB, accountID string, platformID int) ([]model.AccountPlatform, error) {
	accountPlatforms := make([]model.AccountPlatform, 0)
	if err := db.Where("account_id = ? and platform_id = ?", accountID, platformID).Find(&accountPlatforms).Error; err != nil {
		return nil, err
	}

	return accountPlatforms, nil
}

func (d *database) QueryLinks(db *gorm.DB, _type int, identity string, suffixID, pageIndex int) ([]model.Link, error) {
	links := make([]model.Link, 0)
	if err := db.Where(
		"type = ? and identity = ? and suffix_id = ? and page_index = ?",
		_type, identity, suffixID, pageIndex,
	).Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

func (d *database) QueryLinksByTarget(
	db *gorm.DB,
	_type int,
	targetIdentity string,
	targetSuffixID, limit int,
	instance,
	lastInstance string,
) ([]model.Link, error) {
	links := make([]model.Link, 0)
	query := db.Where("type = ? and target_identity = ? and target_suffix_id = ?", _type, targetIdentity, targetSuffixID)

	if limit >= 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&links).Error; err != nil {
		return nil, err
	}

	return links, nil
}

func (d *database) QueryLinkList(db *gorm.DB, _type int, identity string, prefixID, suffixID int) (*model.LinkList, error) {
	linkList := model.LinkList{
		Type:     _type,
		Identity: identity,
		PrefixID: prefixID,
		SuffixID: suffixID,
	}
	if err := db.First(&linkList).Error; err != nil {
		return nil, err
	}

	return &linkList, nil
}

func (d *database) QuerySignature(db *gorm.DB, fileURI string) (*model.Signature, error) {
	signature := model.Signature{}
	if err := db.Where("file_uri = ?", fileURI).First(&signature).Error; err != nil {
		return nil, err
	}

	return &signature, nil
}
