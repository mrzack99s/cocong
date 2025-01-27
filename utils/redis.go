package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func RedisCountKeysByPrefix(ctx context.Context, rdb *redis.Client, prefix string) (int64, error) {
	cursor := uint64(0)
	count := 0

	for {
		var keys []string
		var err error

		keys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return 0, err
		}

		count += len(keys)

		if cursor == 0 {
			break
		}
	}

	return int64(count), nil
}

// type RedisData[T any] struct {
// 	Data T `json:"data"`
// }

func RedisSet[T any](ctx context.Context, rdb *redis.Client, key string, data T, expiration time.Duration) error {

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, key, string(b), expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func RedisGet[T any](ctx context.Context, rdb *redis.Client, key string) (T, error) {

	var result T
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return result, fmt.Errorf("key '%s' does not exist", key)
	} else if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(val), &result); err == nil {
		return result, nil
	}

	return result, nil
}

func RedisGetInsideWildcard[T any](ctx context.Context, rdb *redis.Client, pattern string) (T, error) {

	var result T

	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return result, err
		}

		if len(keys) > 0 {
			value, err := rdb.Get(ctx, keys[0]).Result()
			if err != nil {
				return result, err
			}

			if err := json.Unmarshal([]byte(value), &result); err != nil {
				return result, err
			} else {
				break
			}

		}

		if cursor == 0 {
			return result, errors.New("pattern not found")
		}
	}

	return result, nil
}

func GetKeysByPrefix[T any](ctx context.Context, rdb *redis.Client, prefix string) ([]T, error) {
	var results []T
	cursor := uint64(0)

	for {
		var currentKeys []string
		var err error

		currentKeys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range currentKeys {

			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}

			var data T
			if err := json.Unmarshal([]byte(val), &data); err != nil {
				return nil, err
			}
			results = append(results, data)
		}

		if cursor == 0 {
			break
		}
	}

	return results, nil
}

func RedisGetKeysByPrefixWithOffset[T any](ctx context.Context, rdb *redis.Client, prefix string, offset int, limit int) ([]T, error) {
	var results []T
	cursor := uint64(0)

	totalKeys := 0

	for {
		var currentKeys []string
		var err error

		currentKeys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, err
		}

		totalKeys += len(currentKeys)

		if cursor == 0 {
			break
		}
	}
	page := offset/limit + 1
	start := (page - 1) * limit
	if start < 0 {
		start = 0
	}

	cursor = 0
	currentPage := 0

	for {
		var currentKeys []string
		var err error

		currentKeys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range currentKeys {
			if currentPage >= start && len(results) < limit {
				val, err := rdb.Get(ctx, key).Result()
				if err != nil {
					return nil, err
				}

				var data T
				if err := json.Unmarshal([]byte(val), &data); err != nil {
					return nil, err
				}
				results = append(results, data)
			}
			currentPage++
		}

		if cursor == 0 {
			break
		}
	}

	return results, nil
}

func RedisSearchKeysByPrefix[T any](ctx context.Context, rdb *redis.Client, prefix string) ([]T, error) {
	var results []T
	cursor := uint64(0)

	for {
		var currentKeys []string
		var err error

		currentKeys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range currentKeys {
			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}

			var data T
			if err := json.Unmarshal([]byte(val), &data); err != nil {
				return nil, err
			}
			results = append(results, data)

		}

		if cursor == 0 {
			break
		}
	}

	return results, nil
}

func RedisSearchKeysByContain[T any](ctx context.Context, rdb *redis.Client, prefix, substring string) ([]T, error) {
	var results []T
	cursor := uint64(0)

	for {
		var currentKeys []string
		var err error

		currentKeys, cursor, err = rdb.Scan(ctx, cursor, prefix+"*", 0).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range currentKeys {
			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}

			if strings.Contains(val, substring) {
				var data T
				if err := json.Unmarshal([]byte(val), &data); err != nil {
					return nil, err
				}
				results = append(results, data)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return results, nil
}

func RedisKeyExists(ctx context.Context, rdb *redis.Client, key string) (bool, error) {
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func RedisUpdateTTL(ctx context.Context, rdb *redis.Client, key string, ttl time.Duration) error {
	updated, err := rdb.Expire(ctx, key, ttl).Result()
	if err != nil {
		return err
	}
	if !updated {
		return fmt.Errorf("failed to update TTL: key '%s' does not exist", key)
	}
	return nil
}
