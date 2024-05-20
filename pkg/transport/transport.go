package transport

import "github.com/IBM/sarama"

type Transport struct {
	Name    string
	Brokers []string
	Topic   string

	Config   *sarama.Config
	producer sarama.SyncProducer
}
