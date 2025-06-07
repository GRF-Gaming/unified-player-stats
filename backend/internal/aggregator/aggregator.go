package aggregator

import (
	"backend/internal/models"
	"backend/internal/queue"
	"backend/internal/utils/env_var"
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"log/slog"
	"sync"
)

var (
	aggregatorSingleton *Aggregator
	aggOnce             sync.Once = sync.Once{}
)

type Aggregator struct {
	KafkaProducer *kgo.Client
}

func GetAggregator() *Aggregator {
	aggOnce.Do(func() {

		// Load environment variables
		env := env_var.GetAggVars()

		// Create kafka producer for aggregator for kills
		kfkProd, err := queue.NewKafkaProduceConn(env.AggKafkaSeeds, env.AggTopicKills)
		if err != nil {
			log.Fatal(err)
		}

		// Init agg
		aggregatorSingleton = &Aggregator{
			KafkaProducer: kfkProd,
		}
	})

	return aggregatorSingleton
}

func (a *Aggregator) Serve(port int) {

	h := server.Default(server.WithHostPorts(fmt.Sprintf("0.0.0.0:%d", port)))

	h.POST("/update", handleEvents)

	h.Spin()
}

func (a *Aggregator) CleanUp() {
	slog.Info("Closing Kafka connection in aggregator")
	a.KafkaProducer.Close()
}

func (a *Aggregator) EmitKillEvent(ctx context.Context, killEvent *models.EventKillRecord) error {

	// Serialise kill event
	jsonData, err := sonic.Marshal(killEvent)
	if err != nil {
		slog.Error("Unable to serialise kill event", "event", killEvent, "err", err)
		return err
	}

	// Create kafka record
	record := &kgo.Record{
		Value: jsonData,
	}
	res := a.KafkaProducer.ProduceSync(ctx, record)
	if len(res) != 1 {
		slog.Error("Unable to produce to kafka kill topic")
		return errors.New("unable to produce to kafka")
	}

	if res[0].Err != nil {
		slog.Error("Failed to produce kill record", "record", record, "err", res[0].Err)
		return errors.New("unable to produce to kafka")
	}

	return nil
}

func (a *Aggregator) EmitBatchKillEvent(ctx context.Context, killEvents []*models.EventKillRecord) error {

	// TODO: improve batching submission of records

	for _, e := range killEvents {
		slog.Info("Sending kill event", "kill_event", e)
		if err := a.EmitKillEvent(ctx, e); err != nil {
			return err
		}
	}
	return nil
}
