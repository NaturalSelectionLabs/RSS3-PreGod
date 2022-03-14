package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Asset struct {
	ID      primitive.ObjectID `bson:"_id"`
	Content []string
	Path    string
}
