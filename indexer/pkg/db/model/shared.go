package model

import "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"

type Metadata map[string]interface{}

type ObjectId struct {
	NetworkID constants.NetworkID `json:"network_id" bson:"network_id"`
	Proof     string              `json:"proof" bson:"proof"`
}
