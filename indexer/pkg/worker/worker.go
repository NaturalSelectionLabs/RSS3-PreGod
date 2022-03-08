package worker

import (
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/moralis"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
)

func makeCrawlers() []crawler.Crawler {
	mw := moralis.NewMoralisCrawler()

	return []crawler.Crawler{mw}
}

func ProcessTask(t *Task) error {
	crawlers := makeCrawlers()
	for _, c := range crawlers {
		err := c.Work(t.Identity, t.Network)

		if err != nil {
			panic(err)
		}

		r := c.GetResult()

		for _, item := range r.Items {
			db.InsertItemDoc(item)
		}

		//TODO: save by account: <identity>@<platform>
		db.SetAssets(t.Identity, r.Assets)
		db.AppendNotes(t.Identity, r.Notes)
	}

	return nil
}
