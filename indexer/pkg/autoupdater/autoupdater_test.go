package autoupdater_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/autoupdater"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/crawler"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/indexer/pkg/db"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	_ "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/constants"
	_ "github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack"
)

func init() {
	if err := cache.Setup(); err != nil {
		log.Fatalf("cache.Setup err: %v", err)
	}

	cache.GetRedisClient().FlushDB(context.Background())

	if err := db.Setup(); err != nil {
		log.Fatalf("web.Setup err: %v", err)
	}
}

func Test_Run(t *testing.T) {
	ctx := context.Background()

	// Test Add
	for i := 0; i < 50; i++ {
		autoupdater.AddToRecentVisitQueue(ctx, &crawler.WorkParam{
			Identity:   "henry@nya.one",
			NetworkID:  constants.NetworkIDMisskey,
			PlatformID: constants.PlatformIDMisskey,
			Limit:      100,
			Timestamp:  time.Now().Add(-time.Hour),
		})
	}
	time.Sleep(time.Second)
	ts := time.Now().Unix()
	rdb := cache.GetRedisClient()
	zitems, err := rdb.ZRangeWithScores(ctx, autoupdater.RecentVisitQueue.Key, 0, -1).Result()
	assert.Nil(t, err)

	identity2ts := []int64{}

	for _, zitem := range zitems {
		param := crawler.WorkParam{}
		assert.Nil(t, msgpack.Unmarshal([]byte(zitem.Member.(string)), &param)) // nolint: forcetypeassert // in testing
		identity2ts = append(identity2ts, int64(zitem.Score))
	}

	assert.EqualValues(t, 50, len(identity2ts))

	for _, valueTS := range identity2ts {
		assert.True(t, valueTS < ts)
	}

	// clear old
	items, _, err := rdb.ZScan(ctx, autoupdater.RecentVisitQueue.Key, 0, "*", 0).Result()
	assert.Nil(t, err)

	itemsGot := len(items) / 2

	for index, item := range items {
		if index%2 != 0 { // score
			continue
		}

		rdb.ZAdd(ctx, autoupdater.RecentVisitQueue.Key, []*redis.Z{{
			Score:  10086,
			Member: item,
		}}...)
	}

	count, err := rdb.ZCard(ctx, autoupdater.RecentVisitQueue.Key).Result()
	assert.Nil(t, err)
	assert.EqualValues(t, 50, count)

	autoupdater.RecentVisitQueue.ClearOld(ctx)

	count, err = rdb.ZCard(ctx, autoupdater.RecentVisitQueue.Key).Result()
	assert.Nil(t, err)
	assert.EqualValues(t, 50-itemsGot, count)

	// Iter and Run
	// Remove 35, in most cases there will be 5 left.
	// If left too many, crawl may be too slow
	items, _, _ = rdb.ZScan(ctx, autoupdater.RecentVisitQueue.Key, 0, "*", 35).Result()
	for index, item := range items { // expire 10
		if index%2 != 0 { // score
			continue
		}

		rdb.ZAdd(ctx, autoupdater.RecentVisitQueue.Key, []*redis.Z{{
			Score:  10086,
			Member: item,
		}}...)
	}

	err = autoupdater.RunRecentVisitQueue(ctx)
	assert.Nil(t, err)
	count, err = rdb.ZCard(ctx, autoupdater.RecentVisitQueue.Key).Result()
	assert.EqualValues(t, 50-itemsGot-len(items)/2, count)
	assert.Nil(t, err)
}
