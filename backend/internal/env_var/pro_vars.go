package env_var

import (
	"github.com/caarlos0/env/v11"
	"log"
	"strings"
	"sync"
)

var (
	proVarsOnce               = sync.Once{}
	proVarsSingleton *ProVars = &ProVars{}
)

type ProVars struct {
	ProKafkaSeedsRawString string `env:"PRO_KAFKA_SEEDS,required,notEmpty"`
	ProKafkaSeeds          []string
	ProKafkaGroup          string `env:"PRO_KAFKA_GROUP,required,notEmpty"`
	ProTopicKills          string `env:"PRO_TOPIC_KILLS,required,notEmpty"`
	ProKafkaId             string `env:"PRO_KAFKA_ID,required,notEmpty"`

	ProRedisAddr          string `env:"PRO_REDIS_ADDR,required,notEmpty"`
	ProRedisPort          int    `env:"PRO_REDIS_PORT,required,notEmpty"`
	ProRedisPassword      string `env:"PRO_REDIS_PASSWORD" envDefault:""`
	ProRedisDbNumber      int    `env:"PRO_REDIS_DB_NUMBER,required,notEmpty"`
	ProRedisMaxActiveConn int    `env:"PRO_REDIS_MAX_ACTIVE_CONN" envDefault:"5"`
}

func (p *ProVars) validateAndPopulate() error {

	p.ProKafkaSeeds = strings.Split(p.ProKafkaSeedsRawString, ",")

	return nil
}

func GetProVars() *ProVars {
	proVarsOnce.Do(func() {
		if err := env.Parse(proVarsSingleton); err != nil {
			log.Fatal(err)
		}

		if err := proVarsSingleton.validateAndPopulate(); err != nil {
			log.Fatal(err)
		}
	})

	return proVarsSingleton
}
