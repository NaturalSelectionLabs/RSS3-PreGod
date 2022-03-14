package misskey_test

import (
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/api/misskey"
	"github.com/stretchr/testify/assert"
)

type benchmark struct {
	id        string
	tsp       time.Time
	text      string
	fileCount int
}

// TODO: add more tests maybe
func TestGetNoteList(t *testing.T) {
	t.Parallel()

	var benchmarkList []benchmark

	tsp, _ := time.Parse(time.RFC3339, "2022-03-11T12:51:04.739Z")

	benchmarkList = append(benchmarkList, benchmark{
		id:        "8xpzevfn1o",
		tsp:       tsp,
		text:      "And my car Carrot.<img class=\"media\" src=\"https://file.nya.one/misskey/webpublic-5cd68bdc-941a-4df4-b5d9-1445dc5b88ea.jpg\">",
		fileCount: 1,
	})

	tsp, _ = time.Parse(time.RFC3339, "2022-03-11T12:50:50.604Z")

	benchmarkList = append(benchmarkList, benchmark{
		id:        "8xpzekj01m",
		tsp:       tsp,
		text:      "This is my cat Professor.<img class=\"media\" src=\"https://file.nya.one/misskey/webpublic-e6f17f24-f2b0-42e6-8dcc-8942738a126a.jpg\">",
		fileCount: 1,
	})

	tsp, _ = time.Parse(time.RFC3339, "2022-03-11T12:50:23.428Z")

	benchmarkList = append(benchmarkList, benchmark{
		id:        "8xpzdzk41i",
		tsp:       tsp,
		text:      "Yo, I'm Henry from RSS3.",
		fileCount: 0,
	})

	noteList, err := misskey.GetUserNoteList("henry@nya.one", 100, time.Now())

	assert.Nil(t, err)
	assert.Equal(t, 3, len(noteList))

	for k, node := range noteList {
		assert.Equal(t, benchmarkList[k].id, node.Id)
		assert.Equal(t, benchmarkList[k].text, node.Summary)
		assert.Equal(t, benchmarkList[k].tsp, node.CreatedAt)
		assert.Equal(t, benchmarkList[k].fileCount, len(node.Attachments))
	}

	emptyList, err := misskey.GetUserNoteList("henry@nya.one", 100, time.Now().Add(-time.Hour*24*10))

	assert.Nil(t, err)
	assert.Equal(t, 0, len(emptyList))
}
