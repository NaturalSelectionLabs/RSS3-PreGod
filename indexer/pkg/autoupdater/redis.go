/*
This package is to update recent users' data
*/
package autoupdater

import (
	"context"
	"strconv"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/cache"
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

func (q *RedisRecentQueue) Add(ctx context.Context, item string) error {
	return cache.ZAdd(ctx, q.Key, item, float64(time.Now().Unix()))
}

func (q *RedisRecentQueue) Iter(ctx context.Context, runner func(string) error) error {
	var (
		items  []string
		cursor uint64
		err    error
		result error
	)

	if err = q.clearOld(ctx); err != nil {
		result = multierr.Append(result, err)
	}

	for {
		items, cursor, err = cache.ZScan(ctx, q.Key, cursor, "*", IterCount)
		if err != nil {
			result = multierr.Append(result, err)
		}

		for _, item := range items {
			if err := runner(item); err != nil {
				result = multierr.Append(result, err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return result
}

// Those who have not come for a long time( > Duration) will be deleted
func (q *RedisRecentQueue) clearOld(ctx context.Context) error {
	oldestTime := int(time.Now().Add(-q.Duration).Unix())

	return cache.ZRemRangeByScore(ctx, q.Key, strconv.Itoa(0), strconv.Itoa(oldestTime))
}
