package model

import "github.com/kamva/mgm/v3"

type AccountItemList struct {
	mgm.DefaultModel `bson:",inline"`

	AccountInstance string `json:"account_instance" bson:"account_instance"`

	Profiles []Profile `json:"profiles" bson:"profiles"`

	Assets []ObjectId `json:"assets" bson:"assets"`

	Notes []ObjectId `json:"notes" bson:"notes"`
}
