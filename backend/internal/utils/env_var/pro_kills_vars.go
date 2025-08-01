package env_var

import (
	"github.com/caarlos0/env/v11"
	"log"
	"strings"
	"sync"
)

var (
	proVarsOnce                    = sync.Once{}
	proVarsSingleton *ProKillsVars = &ProKillsVars{}
)

type ProKillsVars struct {
	ProKillsKafkaSeedsRawString string `env:"PRO_KILLS_KAFKA_SEEDS,required,notEmpty"`
	ProKillsKafkaSeeds          []string
	ProKillsKafkaGroup          string `env:"PRO_KILLS_KAFKA_GROUP,required,notEmpty"`
	ProKillsTopicKills          string `env:"PRO_KILLS_TOPIC_KILLS,required,notEmpty"`
	ProKillsKafkaId             string `env:"PRO_KILLS_KAFKA_ID,required,notEmpty"`

	ProKillsRedisAddr          string `env:"PRO_KILLS_REDIS_ADDR,required,notEmpty"`
	ProKillsRedisPort          int    `env:"PRO_KILLS_REDIS_PORT,required,notEmpty"`
	ProKillsRedisPassword      string `env:"PRO_KILLS_REDIS_PASSWORD" envDefault:""`
	ProKillsRedisDbNumber      int    `env:"PRO_KILLS_REDIS_DB_NUMBER,required,notEmpty"`
	ProKillsRedisMaxActiveConn int    `env:"PRO_KILLS_REDIS_MAX_ACTIVE_CONN" envDefault:"5"`

	ProKillsInfluxDbAddr     string `env:"PRO_KILLS_INFLUXDB_ADDR"`
	ProKillsInfluxDbPort     int    `env:"PRO_KILLS_INFLUXDB_PORT"`
	ProKillsInfluxDbUsername string `env:"INFLUXDB_USERNAME"`
	ProKillsInfluxDbPassword string `env:"INFLUXDB_PASSWORD"`
	ProKillsInfluxDbOrg      string `env:"INFLUXDB_ORG"`
	ProKillsInfluxDbBucket   string `env:"INFLUXDB_BUCKET"`
	ProKillsInfluxDbToken    string `env:"INFLUXDB_ADMIN_TOKEN"`
}

func (p *ProKillsVars) validateAndPopulate() error {

	p.ProKillsKafkaSeeds = strings.Split(p.ProKillsKafkaSeedsRawString, ",")

	return nil
}

func GetProKillsVars() *ProKillsVars {
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
