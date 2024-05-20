package transport

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/nenavizhuleto/horizon/protocol"
)

type DetectionTransport Transport

func NewDetectionTransport(name string, brokers []string, topic string) (*DetectionTransport, error) {
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

	return &DetectionTransport{
		Name:     name,
		Brokers:  brokers,
		Topic:    topic,
		Config:   config,
		producer: producer,
	}, nil
}

func (t *DetectionTransport) Send(detection protocol.MotionDetectionMessage) error {
	value, err := json.Marshal(detection)
	if err != nil {
		return err
	}

	_, _, err = t.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     t.Topic,
		Timestamp: detection.Timestamp,
		Value:     sarama.ByteEncoder(value),
	})
	return err
}
