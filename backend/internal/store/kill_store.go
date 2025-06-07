package store

import (
	"backend/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"strconv"
	"time"
)

const (
	totalKillsKeyFmt         = "%s:%s:player:kills:total:%s"          // gameName, serverName, playerId
	totalFriendlyKillsKeyFmt = "%s:%s:player:friendly_kills:total:%s" // gameName, serverName, playerId
	totalDeathsKeyFmt        = "%s:%s:player:deaths:total:%s"         // gameName, serverName, playerId

	rolling48hKillsKeyFmt         = "%s:%s:player:kills:rolling48h:%d:%s"          // gameName, serverName, epochHours, playerId
	rolling48hFriendlyKillsKeyFmt = "%s:%s:player:friendly_kills:rolling48h:%d:%s" // gameName, serverName, epochHours, playerId
	rolling48hDeathsKeyFmt        = "%s:%s:player:deaths:rolling48h:%d:%s"         // gameName, serverName, epochHours, playerId
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

func (k *KillStore) RecordKill(
	ctx context.Context,
	gameName string,
	serverName string,
	time time.Time,
	playerId string,
) error {

	var recordErr error

	// Incr global kill count
	key := fmt.Sprintf(totalKillsKeyFmt, gameName, serverName, playerId)
	if err := k.redisConn.Incr(ctx, key).Err(); err != nil {
		slog.Error("Unable to globally incr total kill count", "player_id", playerId, "err", err)
		recordErr = errors.Join(err)
	}

	// Incr rolling kill count
	eHr := utils.GetEpochHour(time)
	rKey := fmt.Sprintf(rolling48hKillsKeyFmt, gameName, serverName, eHr, playerId)
	hrCount, err := k.redisConn.Incr(ctx, rKey).Result()
	if hrCount == 1 { // first kill of the hour
		slog.Debug("Received first rolling kill record for player, setting key expr", "time", time, "player_id", playerId)
		if expErr := k.redisConn.ExpireAt(ctx, rKey, utils.GetTimeFromEpochHour(eHr+50)).Err(); expErr != nil {
			slog.Error("Unable to set expiry for key", "key", key, "err", expErr)
			recordErr = errors.Join(expErr)
		}
	}
	if err != nil {
		slog.Error("Unable to incr rolling kill count", "player_id", "err", err)
		recordErr = errors.Join(err)
	}

	return recordErr
}

func (k *KillStore) RecordFriendlyKill(ctx context.Context, gameName string, serverName string, time time.Time, playerId string) error {

	var recordErr error

	// Incr global count
	key := fmt.Sprintf(totalFriendlyKillsKeyFmt, gameName, serverName, playerId)
	if err := k.redisConn.Incr(ctx, key).Err(); err != nil {
		slog.Error("Unable to incr total friendly kill count", "player_id", playerId, "err", err)
		recordErr = errors.Join(err)
	}

	// Incr rolling-friendly kill count
	eHr := utils.GetEpochHour(time)
	rKey := fmt.Sprintf(rolling48hFriendlyKillsKeyFmt, gameName, serverName, eHr, playerId)
	hrCount, err := k.redisConn.Incr(ctx, rKey).Result()
	if hrCount == 1 { // first kill of the hour
		slog.Debug("Received first rolling friendly kill record for player, setting key expr", "time", time, "player_id", playerId)
		if expErr := k.redisConn.ExpireAt(ctx, rKey, utils.GetTimeFromEpochHour(eHr+50)).Err(); expErr != nil {
			slog.Error("Unable to set expiry for key", "key", key, "err", expErr)
			recordErr = errors.Join(expErr)
		}
	}
	if err != nil {
		slog.Error("Unable to incr rolling friendly kill count", "player_id", "err", err)
		recordErr = errors.Join(err)
	}

	return recordErr
}

func (k *KillStore) RecordDeath(ctx context.Context, gameName string, serverName string, time time.Time, playerId string) error {
	var recordErr error

	// Incr global count
	key := fmt.Sprintf(totalDeathsKeyFmt, gameName, serverName, playerId)
	if err := k.redisConn.Incr(ctx, key).Err(); err != nil {
		slog.Error("Unable to incr total death count", "player_id", playerId, "err", err)
		recordErr = errors.Join(err)
	}

	// Incr rolling-friendly death
	eHr := utils.GetEpochHour(time)
	rKey := fmt.Sprintf(rolling48hDeathsKeyFmt, gameName, serverName, eHr, playerId)
	hrCount, err := k.redisConn.Incr(ctx, rKey).Result()
	if hrCount == 1 { // first kill of the hour
		slog.Debug("Received first rolling death record for player, setting key expr", "time", time, "player_id", playerId)
		if expErr := k.redisConn.ExpireAt(ctx, rKey, utils.GetTimeFromEpochHour(eHr+50)).Err(); expErr != nil {
			slog.Error("Unable to set expiry for key", "key", key, "err", expErr)
			recordErr = errors.Join(expErr)
		}
	}
	if err != nil {
		slog.Error("Unable to incr rolling death count", "player_id", "err", err)
		recordErr = errors.Join(err)
	}

	return recordErr
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

func (k *KillStore) GetTotalKills(ctx context.Context, gameName string, serverName string, playerId string) (int64, error) {
	key := fmt.Sprintf(totalKillsKeyFmt, gameName, serverName, playerId)
	res, err := k.executeNumericRedisQuery(ctx, key)
	if err != nil {
		slog.Error("Unable to query for total kills")
	}
	return res, nil
}

func (k *KillStore) GetTotalFriendlyKills(ctx context.Context, gameName string, serverName string, playerId string) (int64, error) {
	key := fmt.Sprintf(totalFriendlyKillsKeyFmt, gameName, serverName, playerId)
	res, err := k.executeNumericRedisQuery(ctx, key)
	if err != nil {
		slog.Error("Unable to query for total friendly kills")
	}
	return res, nil
}

func (k *KillStore) GetTotalDeaths(ctx context.Context, gameName string, serverName string, playerId string) (int64, error) {
	key := fmt.Sprintf(totalDeathsKeyFmt, gameName, serverName, playerId)
	res, err := k.executeNumericRedisQuery(ctx, key)
	if err != nil {
		slog.Error("Unable to query for total deaths")
	}
	return res, nil
}
