package model

import (
	"github.com/kamva/mgm/v3"
)

type AccountItemList struct {
	mgm.DefaultModel `bson:",inline"`

	AccountInstance string `json:"account_instance" bson:"account_instance"`

	Assets []ItemId `json:"assets" bson:"assets"`

	Notes []ItemId `json:"notes" bson:"notes"`
}
