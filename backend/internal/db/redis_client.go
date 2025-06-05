package db

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

func NewRedisClient(addr string, port int, password string, db int, maxActiveConns int) (*redis.Client, error) {
	if addr == "" {
		return nil, DbInvalidAddr
	}

	if port < 0 {
		return nil, DbInvalidPort
	}

	if maxActiveConns < 0 {
		return nil, DbInvalidMaxConnSize
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:           fmt.Sprintf("%s:%d", addr, port),
		Password:       password,
		DB:             db,
		MaxActiveConns: maxActiveConns,
	})

	if rdb == nil {
		slog.Error("Unable to initialize Redis client")
		return nil, DbGenericErr
	}

	return rdb, nil
}
