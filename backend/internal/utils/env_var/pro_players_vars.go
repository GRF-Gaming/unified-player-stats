package env_var

import (
	"github.com/caarlos0/env/v11"
	"log"
	"strings"
	"sync"
)

var (
	proPlayersOnce                      = sync.Once{}
	proPlayersSingleton *ProPlayersVars = &ProPlayersVars{}
)

type ProPlayersVars struct {
	ProPlayersCacheTTL int `env:"PRO_PLAYERS_CACHE_TTL" envDefault:"15000"`

	ProPlayersKafkaSeedsRawString string `env:"PRO_PLAYERS_KAFKA_SEEDS,required,notEmpty"`
	ProPlayersKafkaSeeds          []string
	ProPlayersKafkaGroup          string `env:"PRO_PLAYERS_KAFKA_GROUP,required,notEmpty"`
	ProPlayersKafkaTopic          string `env:"PRO_PLAYERS_TOPIC_NAMES,required,notEmpty"`
	ProPlayersKafkaId             string `env:"PRO_PLAYERS_KAFKA_ID,required,notEmpty"`

	ProPlayersRedisAddr          string `env:"PRO_PLAYERS_REDIS_ADDR,required,notEmpty"`
	ProPlayersRedisPort          int    `env:"PRO_PLAYERS_REDIS_PORT,required,notEmpty"`
	ProPlayersRedisPassword      string `env:"PRO_PLAYERS_REDIS_PASSWORD" envDefault:""`
	ProPlayersRedisDbNumber      int    `env:"PRO_PLAYERS_REDIS_DB_NUMBER,required,notEmpty"`
	ProPlayersRedisMaxActiveConn int    `env:"PRO_PLAYERS_REDIS_MAX_ACTIVE_CONN" envDefault:"5"`
}

func (p *ProPlayersVars) validateAndPopulate() error {
	p.ProPlayersKafkaSeeds = strings.Split(p.ProPlayersKafkaSeedsRawString, ",")
	return nil
}

func GetProNameVars() *ProPlayersVars {
	proPlayersOnce.Do(func() {
		if err := env.Parse(proPlayersSingleton); err != nil {
			log.Fatal(err)
		}
		if err := proPlayersSingleton.validateAndPopulate(); err != nil {
			log.Fatal(err)
		}
	})
	return proPlayersSingleton
}
