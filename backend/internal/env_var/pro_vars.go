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
