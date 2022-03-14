package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID      primitive.ObjectID `bson:"_id"`
	Content FileContent
	Path    string
}

type FileContent struct {
	Version        string
	ID             string
	DateCreated    time.Time `bson:"date_created"`
	DateUpdated    time.Time `bson:"date_updated"`
	Signature      string
	Profile        FileContentProfile
	Links          []FileContentLink
	List           []string
	Backlinks      []FileContentBackLink
	Pass           FileContentPass `bson:"_pass"`
	Items          FileContentItems
	Assets         FileContentAssets
	AgentID        string
	AgentSignature string
}

type FileContentProfile struct {
	Name     string
	Avatar   []string
	Bio      string
	Accounts []FileContentProfileAccount
}

type FileContentProfileAccount struct {
	ID        string
	Signature string
	Tags      []string
}

type FileContentLink struct {
	ID   string
	List string
}

type FileContentBackLink struct {
	Auto bool
	ID   string
	List string
}

type FileContentPass struct {
	Assets []FileContentPassAsset
}

type FileContentPassAsset struct {
	ID    string
	Class string
	Order int
}

type FileContentItems struct {
	ListAuto   string
	ListCustom string
}

type FileContentAssets struct {
	ListAuto string
}
