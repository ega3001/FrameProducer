package transport

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/nenavizhuleto/horizon/protocol"
)

type FrameTransport Transport

func NewFrameTransport(name string, brokers []string, topic string) (*FrameTransport, error) {
	config := sarama.NewConfig()
	{
		config.Producer.RequiredAcks = sarama.NoResponse
		config.Producer.Return.Successes = true
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &FrameTransport{
		Name:     name,
		Brokers:  brokers,
		Topic:    topic,
		Config:   config,
		producer: producer,
	}, nil
}

func (t *FrameTransport) Send(message protocol.FrameMessage) error {
	value, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, _, err = t.producer.SendMessage(&sarama.ProducerMessage{
		Topic: t.Topic,
		Value: sarama.ByteEncoder(value),
	})
	return err
}
