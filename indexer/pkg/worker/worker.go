package worker

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
)

// TODO: move this to indexer/pkg/api/moralis/worker.go
func ProcessTask(t *Task) error {
	mw := moralis.NewMoralisCrawler()
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
