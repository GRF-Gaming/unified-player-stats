package env_var

import (
	"github.com/caarlos0/env/v11"
	"log"
	"strings"
	"sync"
)

var (
	aggVarsOnce               = sync.Once{}
	aggVarsSingleton *AggVars = &AggVars{}
)

type AggVars struct {
	AggKafkaSeedsRawString string `env:"AGG_KAFKA_SEEDS,required,notEmpty"` // Comma seperated strings
	AggKafkaSeeds          []string
	AggTopicKills          string `env:"AGG_TOPIC_KILLS,required,notEmpty"`

	AggPort int `env:"AGG_PORT,required,notEmpty"`
}

// validateAndPopulate handles the parsing and validation of user defined envs
func (a *AggVars) validateAndPopulate() error {

	a.AggKafkaSeeds = strings.Split(a.AggKafkaSeedsRawString, ",")

	return nil
}

func GetAggVars() *AggVars {
	aggVarsOnce.Do(func() {
		if err := env.Parse(aggVarsSingleton); err != nil {
			log.Fatal(err)
		}
		if err := aggVarsSingleton.validateAndPopulate(); err != nil {
			log.Fatal(err)
		}
	})
	return aggVarsSingleton
}
