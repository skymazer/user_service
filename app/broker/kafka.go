package broker

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/skymazer/user_service/loggerfx"
)

type Kafka struct {
	Conn *kafka.Conn
	log  *loggerfx.Logger
}

const (
	host      = "kafka"
	port      = "9092"
	partition = 0
)

func New(log *loggerfx.Logger, topic string) (*Kafka, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:%s", host, port), topic, partition)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial leader:")
	}

	return &Kafka{conn, log}, nil

}

func (k *Kafka) WriteLog(message []byte) error {
	_, err := k.Conn.WriteMessages(
		kafka.Message{
			Key:   []byte("SomeKey"),
			Value: []byte(message)},
	)

	return errors.Wrap(err, "failed to write message:")
}
