package db

import (
	"fmt"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/rss3uri"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Setup() error {
	err := mgm.SetDefaultConfig(nil, config.Config.Mongo.DB, options.Client().ApplyURI(config.Config.Mongo.URI))
	// TODO: Create Indexes
	// Reference: https://github.com/Kamva/mgm/issues/35
	return err
}

// SetAssets refresh users' all assets by network
func SetAssets(instance rss3uri.Instance, assets []*model.ItemId, refreshBy constants.NetworkID) {
	mgm.Coll(&model.AccountItemList{}).FindOneAndUpdate(
		mgm.Ctx(), bson.M{"account_instance": instance},
		bson.M{"$pull": bson.M{"assets.network_id": refreshBy}},
	)
	mgm.Coll(&model.AccountItemList{}).FindOneAndUpdate(
		mgm.Ctx(), bson.M{"account_instance": instance},
		bson.M{"$addToSet": bson.M{"assets": bson.M{"$each": assets}}},
		options.FindOneAndUpdate().SetUpsert(true),
	)
}

// AppendNotes only append users' new notes(duplicated notes are omitted.)
func AppendNotes(instance rss3uri.Instance, notes []*model.ItemId) {
	// If the value is a document, MongoDB determines that the document is a duplicate if an existing
	// document in the array matches the to-be-added document exactly; i.e. the existing document has
	// the exact same fields and values and the fields are in the same order. As such, field order
	// matters and you cannot specify that MongoDB compare only a subset of the fields in the document
	// to determine whether the document is a duplicate of an existing array element.
	// If we use ODM, the order is the same. So we do not need to worry
	mgm.Coll(&model.AccountItemList{}).FindOneAndUpdate(
		mgm.Ctx(),
		bson.M{"account_instance": instance.String()}, bson.M{"$addToSet": bson.M{"notes": bson.M{"$each": notes}}},
		options.FindOneAndUpdate().SetUpsert(true),
	)
}

func Exists(i rss3uri.Instance) (bool, error) {
	n, err := mgm.Coll(&model.AccountItemList{}).CountDocuments(
		mgm.Ctx(),
		bson.M{"account_instance": i.String()},
	)
	if err != nil {
		return false, err
	}

	if n == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func GetAccountInstance(instance rss3uri.Instance) (*model.AccountItemList, error) {
	r := &model.AccountItemList{}
	err := mgm.Coll(&model.AccountItemList{}).FindOne(
		mgm.Ctx(),
		bson.M{"account_instance": instance.String()},
	).Decode(r)

	if err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func GetAccountItems(instance rss3uri.Instance, t constants.ItemType) (*[]model.Item, error) {
	r, err := GetAccountInstance(instance)
	if err != nil {
		return nil, err
	}

	idList := []model.ItemId{}
	if t == constants.ItemTypeAsset {
		idList = r.Assets
	} else if t == constants.ItemTypeNote {
		idList = r.Notes
	} else {
		return nil, fmt.Errorf("unsupported instance query")
	}

	if idList == nil {
		return nil, nil
	}

	results, err := GetItems(&idList)
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func InsertItems(items []*model.Item, networkID constants.NetworkID) (*mongo.BulkWriteResult, error) {
	models := []mongo.WriteModel{}

	for _, item := range items {
		models = append(models, mongo.NewReplaceOneModel().SetFilter(
			bson.M{"item_id.network_id": item.ItemId.NetworkID,
				"item_id.proof": item.ItemId.Proof}).SetReplacement(item).SetUpsert(true))
	}

	return mgm.Coll(&model.Item{}).BulkWrite(mgm.Ctx(), models)
}

func GetItem(key *model.ItemId) (*model.Item, error) {
	r := &model.Item{}
	err := mgm.Coll(&model.Item{}).FindOne(
		mgm.Ctx(),
		bson.M{"item_id": key},
	).Decode(r)

	if err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func GetItems(key *[]model.ItemId) (*[]model.Item, error) {
	results := &[]model.Item{}
	err := mgm.Coll(&model.Item{}).SimpleFind(
		results,
		bson.M{"item_id": bson.M{"$in": key}},
	)

	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}
