package worker

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawlers"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
)

func ProcessTask(t *Task) error {
	mw := crawlers.NewMoralisCrawler()
	err := mw.Work(t.Identity, t.Network)
	if err != nil {
		panic(err)
	}
	r := mw.GetResult()

	for _, item := range r.Items {
		db.InsertItemDoc(item)
	}
	for _, object := range r.Objects {
		db.InsertObjectDoc(object)
	}
	//TODO: save by account: <identity>@<platform>
	db.SetAssets(t.Identity, r.Assets)
	db.AppendNotes(t.Identity, r.Notes)

	return nil
}
