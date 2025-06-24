package store

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"log/slog"
	"time"
)

const (
	pointRecordKill         string = "kill"
	pointRecordFriendlyKill string = "ff_kill"
	pointRecordDeath        string = "death"
)

type EventStore struct {
	client   influxdb2.Client
	writeApi api.WriteAPI
}

func NewEventStore(
	addr string,
	port int,
	token string,
	org string,
	bucket string,
) *EventStore {

	slog.Info("[NewEventStore] Created a new event store")

	client := influxdb2.NewClient(fmt.Sprintf("%s:%d", addr, port), token)
	wApi := client.WriteAPI(org, bucket)

	return &EventStore{
		client:   client,
		writeApi: wApi,
	}
}

func (e *EventStore) CleanUp() {
	slog.Info("[EventStore][CleanUp] Closing influxdb client")
	e.writeApi.Flush()
	e.client.Close()
}

func (e *EventStore) RecordKill(
	ctx context.Context,
	gameName string,
	serverName string,
	time time.Time,
	playerId string,
) error {

	point := write.NewPoint(
		pointRecordKill,
		map[string]string{
			"game_name":   gameName,
			"server_name": serverName,
			"player_id":   playerId,
		},
		map[string]interface{}{"count": 1},
		time,
	)

	e.writeApi.WritePoint(point)

	return nil
}

func (e *EventStore) RecordFriendlyKillWithVictim(
	ctx context.Context,
	gameName string,
	serverName string,
	time time.Time,
	playerId string,
	victimPlayerId string,
) error {

	point := write.NewPoint(
		pointRecordFriendlyKill,
		map[string]string{
			"game_name":        gameName,
			"server_name":      serverName,
			"player_id":        playerId,
			"victim_player_id": victimPlayerId,
		},
		map[string]interface{}{"count": 1},
		time,
	)

	e.writeApi.WritePoint(point)

	return nil
}

func (e *EventStore) RecordDeath(
	ctx context.Context,
	gameName string,
	serverName string,
	time time.Time,
	playerId string,
) error {

	point := write.NewPoint(
		pointRecordDeath,
		map[string]string{
			"game_name":   gameName,
			"server_name": serverName,
			"player_id":   playerId,
		},
		map[string]interface{}{"count": 1},
		time,
	)

	e.writeApi.WritePoint(point)

	return nil
}
