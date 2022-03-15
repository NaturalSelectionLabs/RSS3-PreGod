package cache

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/NaturalSelectionLabs/RSS3-PreGod/shared/pkg/config"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

var (
	jsoni = jsoniter.ConfigCompatibleWithStandardLibrary
	rdb   *redis.Client
)

func Setup() error {
	ctx := context.Background()

	var tlsConfig *tls.Config
	if config.Config.Redis.TLSEnabled { // In AWS if a Redis has a password, it must use TLS
		tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:      config.Config.Redis.Addr,
		Password:  config.Config.Redis.Password,
		DB:        config.Config.Redis.DB,
		TLSConfig: tlsConfig,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return err
	}

	return nil
}

func Get(ctx context.Context, key string, data interface{}) error {
	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil || err != nil {
		return err
	}

	return jsoni.Unmarshal([]byte(value), &data)
}

func Set(ctx context.Context, key string, data interface{}, expire time.Duration) error {
	value, err := jsoni.Marshal(data)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, value, expire).Err()
}

func Exists(ctx context.Context, key string) (bool, error) {
	exist, err := rdb.Exists(ctx, key).Result()

	if err != nil {
		return false, err
	}

	return exist == 1, nil
}

func ZAdd(ctx context.Context, key string, data interface{}, score float64) error {
	value, err := jsoni.Marshal(data)
	if err != nil {
		return err
	}

	if err = rdb.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: value,
	}).Err(); err != nil {
		return err
	}

	return nil
}

// zrange in pre-node
func ZRevRange(ctx context.Context, key string, min string, max string, offset int64, count int64) ([]interface{}, error) {
	res := rdb.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	})

	if res == nil {
		return nil, res.Err()
	}

	result := make([]interface{}, len(res.Val()))

	for i, val := range res.Val() {
		err := jsoni.Unmarshal([]byte(val), &result[i])

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// zrangeWithScore in pre-node
func ZRevRangeWithScore(ctx context.Context, key string, min string, max string, offset int64, count int64) ([]interface{}, error) {
	// go-redis will check for Offset and Count
	res := rdb.ZRevRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  count,
	})

	if res == nil {
		return nil, res.Err()
	}

	result := make([]interface{}, len(res.Val()))

	for i, v := range res.Val() {
		val, ok := v.Member.(string)

		if ok {
			err := jsoni.Unmarshal([]byte(val), &result[i])

			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("Unable to convert to string")
		}
	}

	return result, nil
}

// never used in prenode
func ZRem(ctx context.Context, key string, data interface{}, score float64) error {
	value, err := jsoni.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := rdb.ZRem(ctx, key, value).Result(); err != nil {
		return err
	}

	return nil
}

func SRem(ctx context.Context, key string, data interface{}) error {
	value, err := jsoni.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := rdb.SRem(ctx, key, value).Result(); err != nil {
		return err
	}

	return nil
}

func SAdd(ctx context.Context, key string, data interface{}) error {
	value, err := jsoni.Marshal(data)
	if err != nil {
		return err
	}

	if _, err := rdb.SAdd(ctx, key, value).Result(); err != nil {
		return err
	}

	return nil
}

func SGet(ctx context.Context, key string, data interface{}) ([]interface{}, error) {
	res := rdb.SMembers(ctx, key)

	if res == nil {
		return nil, res.Err()
	}

	result := make([]interface{}, len(res.Val()))

	for i, val := range res.Val() {
		if err := jsoni.Unmarshal([]byte(val), &result[i]); err != nil {
			return nil, err
		}
	}

	return result, nil
}
