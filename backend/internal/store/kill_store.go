package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
)

const (
	totalKillsKeyFmt         = "player:kills:total:%s"
	totalFriendlyKillsKeyFmt = "player:friendly_kills:total:%s"
	totalDeathsKeyFmt        = "player:deaths:total:%s"
)

type KillStore struct {
	redisConn *redis.Client
}

func NewKillStore(rdb *redis.Client) *KillStore {
	return &KillStore{
		redisConn: rdb,
	}
}
func (k *KillStore) Close() {
	if err := k.redisConn.Close(); err != nil {
		slog.Error("Unable to close KillStore.redisConn", "err", err)
	}
}

func (k *KillStore) RecordKill(ctx context.Context, playerId string) error {
	key := fmt.Sprintf(totalKillsKeyFmt, playerId)
	return k.redisConn.Incr(ctx, key).Err()
}

func (k *KillStore) RecordFriendlyKill(ctx context.Context, playerId string) error {
	key := fmt.Sprintf(totalFriendlyKillsKeyFmt, playerId)
	return k.redisConn.Incr(ctx, key).Err()
}

func (k *KillStore) RecordDeath(ctx context.Context, playerId string) error {
	key := fmt.Sprintf(totalDeathsKeyFmt, playerId)
	return k.redisConn.Incr(ctx, key).Err()
}

func (k *KillStore) executeNumericRedisQuery(ctx context.Context, key string) (int64, error) {
	val, err := k.redisConn.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil // Player does not have any kills yet
	}
	if err != nil {
		slog.Error("Unexpected error querying redis", "key", key, "err", err)
		return 0, err
	}

	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		slog.Error(
			"Parse error, unable to parse val into i64 in KillStore.executeNumericRedisQuery",
			"key",
			key,
			"val",
			val,
			"err",
			err,
		)
		return 0, err
	}

	return res, nil
}

func (k *KillStore) GetTotalKills(ctx context.Context, playerId string) (int64, error) {
	key := fmt.Sprintf(totalKillsKeyFmt, playerId)
	res, err := k.executeNumericRedisQuery(ctx, key)
	if err != nil {
		slog.Error("Unable to query for total kills")
	}
	return res, nil
}

func (k *KillStore) GetTotalFriendlyKills(ctx context.Context, playerId string) (int64, error) {
	key := fmt.Sprintf(totalFriendlyKillsKeyFmt, playerId)
	res, err := k.executeNumericRedisQuery(ctx, key)
	if err != nil {
		slog.Error("Unable to query for total friendly kills")
	}
	return res, nil
}

func (k *KillStore) GetTotalDeaths(ctx context.Context, playerId string) (int64, error) {
	key := fmt.Sprintf(totalDeathsKeyFmt, playerId)
	res, err := k.executeNumericRedisQuery(ctx, key)
	if err != nil {
		slog.Error("Unable to query for total deaths")
	}
	return res, nil
}
