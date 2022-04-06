package model

import (
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	mgm.DefaultModel `bson:",inline"`

	ProfileId         ObjectId          `json:"profile_id" bson:"profile_id"`
	Metadata          map[string]string `json:"metadata" bson:"metadata"`
	Name              string            `json:"name" bson:"name"`
	Bio               string            `json:"bio" bson:"bio"`
	Avatars           []string          `json:"avatars" bson:"avatars"`
	Attachments       []Attachment      `json:"attachments" bson:"attachments"`
	ConnectedAccounts []string          `json:"connected_accounts" bson:"connected_accounts"`
}

func NewProfile(networkId constants.NetworkID, proof string, metadata map[string]string,
	name string, bio string, avatars []string,
	attachments []Attachment, connectedAccounts []string) *Profile {
	return &Profile{
		DefaultModel: mgm.DefaultModel{
			IDField: mgm.IDField{
				ID: primitive.NewObjectID(),
			},
			DateFields: mgm.DateFields{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		ProfileId: ObjectId{
			NetworkID: networkId,
			Proof:     proof,
			SyncAt:    time.Now(),
		},
		Metadata:          metadata,
		Name:              name,
		Bio:               bio,
		Avatars:           avatars,
		Attachments:       attachments,
		ConnectedAccounts: connectedAccounts,
	}
}
