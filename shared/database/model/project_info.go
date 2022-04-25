package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/database/common"

type ProjectInfo struct {
	Active          bool   `gorm:"column:active"`
	Id              int    `gorm:"column:id"`
	Title           string `gorm:"column:title"`
	Slug            string `gorm:"column:slug"`
	Description     string `gorm:"column:description"`
	ReferUrl        string `gorm:"column:reference_url"`
	Logo            string `gorm:"column:logo"`
	AdminAddress    string `gorm:"column:admin_address"`
	TokenAddress    string `gorm:"column:token_address"`
	TokenSymbol     string `gorm:"column:token_symbol"`
	ContractAddress string `gorm:"column:contract_address"`

	common.Table
}

func (ProjectInfo) TableName() string {
	return "reptile-gitcoin.data"
}
