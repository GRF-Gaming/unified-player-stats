package queue

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"strings"
	"time"
)

func NewKafkaProduceConn(
	seeds []string,
	topic string,
) (*kgo.Client, error) {

	client, err := kgo.NewClient(
		kgo.SeedBrokers(strings.Join(seeds, ",")),
		kgo.DefaultProduceTopic(topic),
		kgo.RecordDeliveryTimeout(10*time.Second),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.AllowAutoTopicCreation(),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewKafkaConsumeConn(
	seeds []string,
	group string,
	topic string,
) (*kgo.Client, error) {

	client, err := kgo.NewClient(
		kgo.SeedBrokers(strings.Join(seeds, ",")),
		kgo.AllowAutoTopicCreation(),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.RecordDeliveryTimeout(10*time.Second),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
