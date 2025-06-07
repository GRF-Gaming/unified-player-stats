package process_players

import (
	"backend/internal/db"
	"backend/internal/queue"
	"backend/internal/store"
	"backend/internal/utils"
	"backend/internal/utils/env_var"
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
	"log"
	"sync"
	"time"
)

var (
	processPlayerSingleton *Processor
	processPlayerOnce      sync.Once = sync.Once{}
)

type Processor struct {
	ctx           context.Context
	cancelFunc    context.CancelFunc
	kafkaConsumer *kgo.Client
	nameStore     *store.NameStore
	nameCache     *utils.TTLCache[string, string]
}

func GetProcessor() *Processor {
	processPlayerOnce.Do(func() {

		e := env_var.GetProNameVars()

		ctx, cancelFunc := context.WithCancel(context.Background())

		kCon, err := queue.NewKafkaConsumeConn(
			e.ProPlayersKafkaId,
			e.ProPlayersKafkaSeeds,
			e.ProPlayersKafkaGroup,
			e.ProPlayersKafkaTopic,
		)
		if err != nil {
			log.Fatal(err)
		}

		rCon, err := db.NewRedisClient(
			e.ProPlayersRedisAddr,
			e.ProPlayersRedisPort,
			e.ProPlayersRedisPassword,
			e.ProPlayersRedisDbNumber,
			e.ProPlayersRedisMaxActiveConn,
		)
		if err != nil {
			log.Fatal(err)
		}

		processPlayerSingleton = &Processor{
			ctx:           ctx,
			cancelFunc:    cancelFunc,
			kafkaConsumer: kCon,
			nameStore:     store.NewNameStore(ctx, rCon),
			nameCache:     utils.NewTTLCache[string, string](time.Duration(e.ProPlayersCacheTTL) * time.Millisecond),
		}

	})

	return processPlayerSingleton
}

func (p *Processor) CleanUp() {
	p.cancelFunc()
	p.kafkaConsumer.Close()
	p.nameStore.CleanUp()
}
