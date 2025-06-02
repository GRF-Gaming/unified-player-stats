package queue

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"strings"
)

func NewKafkaProduceConn(
	seeds []string,
	topic string,
) (*kgo.Client, error) {

	client, err := kgo.NewClient(
		kgo.SeedBrokers(strings.Join(seeds, ",")),
		kgo.DefaultProduceTopic(topic),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewKafkaConsumeConn(
	seeds []string,
	group string,
	topic []string,
) (*kgo.Client, error) {

	client, err := kgo.NewClient(
		kgo.SeedBrokers(strings.Join(seeds, ",")),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic...),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
