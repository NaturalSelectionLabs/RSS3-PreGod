package db

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db/model"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
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
func SetAssets(instance string, assets []*model.ItemId, refreshBy constants.NetworkID) {
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
func AppendNotes(instance string, notes []*model.ItemId) {
	// If the value is a document, MongoDB determines that the document is a duplicate if an existing
	// document in the array matches the to-be-added document exactly; i.e. the existing document has
	// the exact same fields and values and the fields are in the same order. As such, field order
	// matters and you cannot specify that MongoDB compare only a subset of the fields in the document
	// to determine whether the document is a duplicate of an existing array element.
	// If we use ODM, the order is the same. So we do not need to worry
	mgm.Coll(&model.AccountItemList{}).FindOneAndUpdate(
		mgm.Ctx(),
		bson.M{"account_instance": instance}, bson.M{"$addToSet": bson.M{"notes": bson.M{"$each": notes}}},
		options.FindOneAndUpdate().SetUpsert(true),
	)
}

// TODO: getter

func InsertItemDoc(item *model.Item) *mongo.SingleResult {
	return mgm.Coll(&model.Item{}).FindOneAndReplace(
		mgm.Ctx(),
		bson.M{"item_id.network_id": item.ItemId.NetworkId, "item_id.proof": item.ItemId.Proof},
		item,
		options.FindOneAndReplace().SetUpsert(true),
	)
}

// TODO: getter
