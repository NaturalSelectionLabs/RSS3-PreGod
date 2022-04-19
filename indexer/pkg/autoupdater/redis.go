/*
This package is to update recent users' data
*/
package autoupdater

import (
	"context"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/logger"
	"github.com/go-redis/redis/v8"
	"go.uber.org/multierr"
)

const (
	RecentVisitKey      = "index.item.recent.visit"
	RecentVisitDuration = 24 * time.Hour
	IterCount           = 10
)

// redis zset
type RedisRecentQueue struct {
	Key      string
	Duration time.Duration // (`now` - `recent_visit_time`) should be less than Duration
}

var RecentVisitQueue = &RedisRecentQueue{Key: RecentVisitKey, Duration: RecentVisitDuration}

func (q *RedisRecentQueue) Add(ctx context.Context, item []byte) error {
	return cache.GetRedisClient().ZAdd(ctx, q.Key, []*redis.Z{
		{
			Score:  float64(time.Now().Unix()),
			Member: item,
		},
	}...).Err()
}

func (q *RedisRecentQueue) Iter(ctx context.Context, runner func(string) error) error {
	var (
		items  []string
		cursor uint64
		err    error
		result error
	)

	if err = q.ClearOld(ctx); err != nil {
		result = multierr.Append(result, err)
	}

	runCount := 0

	for {
		items, cursor, err = cache.ZScan(ctx, q.Key, cursor, "*", IterCount)
		if err != nil {
			result = multierr.Append(result, err)
		}

		for index, item := range items {
			if index%2 != 0 { // score
				continue
			}

			if err := runner(item); err != nil {
				result = multierr.Append(result, err)
			}

			runCount += 1
		}

		if cursor == 0 {
			break
		}
	}

	logger.Infof("RedisRecentQueue: %d jobs finished", runCount)

	return result
}

// Those who have not come for a long time( > Duration) will be deleted
func (q *RedisRecentQueue) ClearOld(ctx context.Context) error {
	{
		logger.Infof("before clear: %d members in %s", cache.GetRedisClient().ZCard(ctx, q.Key).Val(), q.Key)
	}

	oldestTime := int(time.Now().Add(-q.Duration).Unix())

	num, err := cache.ZRemRangeByScore(ctx, q.Key, strconv.Itoa(0), strconv.Itoa(oldestTime))
	{
		logger.Infof("%d members removed in %s", num, q.Key)
		logger.Infof("after clear: %d members in %s", cache.GetRedisClient().ZCard(ctx, q.Key).Val(), q.Key)
	}

	return err
}
