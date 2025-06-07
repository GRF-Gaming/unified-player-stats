package store

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type NameStore struct {
	ctx     context.Context
	rClient *redis.Client
}

func NewNameStore(ctx context.Context, rClient *redis.Client) *NameStore {
	return &NameStore{
		ctx:     ctx,
		rClient: rClient,
	}
}

func (n *NameStore) CleanUp() {
	if err := n.rClient.Close(); err != nil {
		slog.Error("Unable to close redis connection in NameStore", "err", err)
	}
}

func (k *KillStore) RecordUidNamePair(uid string, playerName string) error {

	return nil
}

func (k *KillStore) GetLatestName(uid string) (string, bool) {

	return "", false
}
