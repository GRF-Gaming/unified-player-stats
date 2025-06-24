package store

import (
	"context"
	"time"
)

type KillStoreInterface interface {
	RecordKillInterface
	RecordFriendlyKillInterface
	RecordDeathInterface
}

type EventStoreInterface interface {
	RecordKillInterface
	RecordFriendlyKillWithVictimInterface
	RecordDeathInterface
}

type RecordKillInterface interface {
	RecordKill(
		ctx context.Context,
		gameName string,
		serverName string,
		time time.Time,
		playerId string,
	) error
}

type RecordFriendlyKillInterface interface {
	RecordFriendlyKill(
		ctx context.Context,
		gameName string,
		serverName string,
		time time.Time,
		playerId string,
	) error
}

type RecordFriendlyKillWithVictimInterface interface {
	RecordFriendlyKillWithVictim(
		ctx context.Context,
		gameName string,
		serverName string,
		time time.Time,
		playerId string,
		victimPlayerId string,
	) error
}

type RecordDeathInterface interface {
	RecordDeath(
		ctx context.Context,
		gameName string,
		serverName string,
		time time.Time,
		playerId string,
	) error
}
