package processor_kills

import (
	"backend/internal/env_var"
	"backend/internal/queue"
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
	"log/slog"
	"sync"
)

var (
	processorSingleton *Processor
	processorOnce      sync.Once = sync.Once{}
)

type Processor struct {
	ctx           context.Context
	cancelFunc    context.CancelFunc
	kafkaConsumer *kgo.Client
}

func GetProcessor() *Processor {
	processorOnce.Do(func() {

		e := env_var.GetProVars()
		consumer, err := queue.NewKafkaConsumeConn(e.ProKafkaId, e.ProKafkaSeeds, e.ProKafkaGroup, e.ProTopicKills)
		if err != nil {
			slog.Error("Unable to create kafka consumer client for kill processor")
		}

		ctx, cancel := context.WithCancel(context.Background())

		processorSingleton = &Processor{
			ctx:           ctx,
			cancelFunc:    cancel,
			kafkaConsumer: consumer,
		}
	})

	return processorSingleton
}

func (p *Processor) Spin() {

	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			fetches := p.kafkaConsumer.PollFetches(p.ctx)
			if fetches.IsClientClosed() {
				slog.Error("Kafka client has been closed, exiting")
				return
			}

			// Handle errors
			if err := fetches.Errors(); len(err) > 0 {
				if p.ctx.Err() != nil {
					// context canceled
					return
				}

				envPro := env_var.GetProVars()

				for _, e := range err {
					slog.Error("Unable to fetch from Kafka", "seeds", envPro.ProKafkaSeeds, "topic", envPro.ProTopicKills, "group", envPro.ProKafkaGroup, "err", e.Err)
				}
			}

			// Process fetched records
			fetches.EachRecord(func(r *kgo.Record) {
				// TODO process record
				slog.Info("Record", "val", string(r.Value))
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
}
