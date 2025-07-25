package processor_kills

import (
	"backend/internal/db"
	"backend/internal/models"
	"backend/internal/queue"
	"backend/internal/store"
	"backend/internal/utils/env_var"
	"backend/internal/utils/pools"
	"context"
	"github.com/bytedance/sonic"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"log/slog"
	"sync"
)

var (
	processorSingleton *Processor
	processorOnce      = sync.Once{}
)

type Processor struct {
	ctx           context.Context
	cancelFunc    context.CancelFunc
	kafkaConsumer *kgo.Client
	killStore     *store.KillStore
	eventStore    *store.EventStore
}

func GetProcessor() *Processor {
	processorOnce.Do(func() {

		e := env_var.GetProKillsVars()

		consumer, err := queue.NewKafkaConsumeConn(
			e.ProKillsKafkaId,
			e.ProKillsKafkaSeeds,
			e.ProKillsKafkaGroup,
			e.ProKillsTopicKills,
		)
		if err != nil {
			slog.Error("[GetProcessor] Unable to create kafka consumer client for kill processor")
			log.Fatal("[GetProcessor] Unable to create Processor")
		}

		rClient, err := db.NewRedisClient(
			e.ProKillsRedisAddr,
			e.ProKillsRedisPort,
			e.ProKillsRedisPassword,
			e.ProKillsRedisDbNumber,
			e.ProKillsRedisMaxActiveConn,
		)
		if err != nil {
			slog.Error("[GetProcessor] Unable to create redis conn client for kill processor")
			log.Fatal("[GetProcessor] Unable to create Processor")
		}

		eventStore := store.NewEventStore(
			e.ProKillsInfluxDbAddr,
			e.ProKillsInfluxDbPort,
			e.ProKillsInfluxDbToken,
			e.ProKillsInfluxDbOrg,
			e.ProKillsInfluxDbBucket,
		)

		ctx, cancel := context.WithCancel(context.Background())

		processorSingleton = &Processor{
			ctx:           ctx,
			cancelFunc:    cancel,
			kafkaConsumer: consumer,
			killStore:     store.NewKillStore(rClient),
			eventStore:    eventStore,
		}
	})

	return processorSingleton
}

func (p *Processor) Spin() {

	envPro := env_var.GetProKillsVars()

	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			fetches := p.kafkaConsumer.PollFetches(p.ctx)
			if fetches.IsClientClosed() {
				slog.Error("[Processor][Spin] Kafka client has been closed, exiting")
				return
			}

			// Handle errors
			if err := fetches.Errors(); len(err) > 0 {
				if p.ctx.Err() != nil {
					// context canceled
					return
				}

				for _, e := range err {
					slog.Error("Unable to fetch from Kafka", "seeds", envPro.ProKillsKafkaSeeds, "topic", envPro.ProKillsTopicKills, "group", envPro.ProKillsKafkaGroup, "err", e.Err)
				}
			}

			// Process fetched records
			fetches.EachRecord(func(r *kgo.Record) {
				killRecord := pools.EventKillRecordPool.Get().(*models.EventKillRecord)
				defer pools.EventKillRecordPool.Put(killRecord)
				if err := sonic.Unmarshal(r.Value, killRecord); err != nil {
					slog.Error("Unable to unmarshal kill record, dropping", "record", string(r.Value))
					p.kafkaConsumer.MarkCommitRecords(r)
					return
				}
				if err := p.ProcessKillRecord(killRecord); err != nil {
					slog.Error("Unable to process record", "event", string(r.Value))
					return
				}
				p.kafkaConsumer.MarkCommitRecords(r)
			})

			// Commit offset
			if err := p.kafkaConsumer.CommitUncommittedOffsets(p.ctx); err != nil {
				slog.Error("Unable to commit offset", "err", err)
			}
			p.kafkaConsumer.AllowRebalance()
		}
	}
}

func (p *Processor) CleanUp() {
	p.cancelFunc()
	p.kafkaConsumer.Close()
	p.killStore.Close()
}

func (p *Processor) ProcessKillRecord(r *models.EventKillRecord) error {

	// Handle updating of KillStore
	if err := p.killStore.RecordKill(p.ctx, r.GameName, r.ServerName, r.Time, r.PlayerId); err != nil {
		slog.Error("[Processor][ProcessKillRecord] Unable to record kill in kill store", "err", err)
	}
	if err := p.eventStore.RecordKill(p.ctx, r.GameName, r.ServerName, r.Time, r.PlayerId); err != nil {
		slog.Error("[Processor][ProcessKillRecord] Unable to record kill in event store", "err", err)
	}

	// Handle events related to friendly kills
	if r.IsFriendly {
		if err := p.killStore.RecordFriendlyKill(p.ctx, r.GameName, r.ServerName, r.Time, r.PlayerId); err != nil {
			slog.Error("[Processor][ProcessKillRecord] Unable to record friendly fire kill in kill store", "err", err)
		}
		if err := p.eventStore.RecordFriendlyKillWithVictim(p.ctx, r.GameName, r.ServerName, r.Time, r.PlayerId, r.KilledEntityId); err != nil {
			slog.Error("[Processor][ProcessKillRecord] Unable to record friendly fire kill in event store", "err", err)
		}
	}

	// Handle deaths of a player
	if err := p.killStore.RecordDeath(p.ctx, r.GameName, r.ServerName, r.Time, r.KilledEntityId); err != nil {
		slog.Error("[Processor][ProcessKillRecord] Unable to record player death in kill store", "err", err)
	}
	if err := p.eventStore.RecordDeath(p.ctx, r.GameName, r.ServerName, r.Time, r.KilledEntityId); err != nil {
		slog.Error("[Processor][ProcessKillRecord] Unable to record player death in event store", "err", err)
	}

	return nil

}
