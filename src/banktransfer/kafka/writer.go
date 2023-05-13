package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/turngeek/myaktion-go-2023/src/banktransfer/grpc/banktransfer"
)

const (
	groupID = "banktransfers-watchers"
)

var connect = os.Getenv("KAFKA_CONNECT")

type TransactionWriter interface {
	Close() error
	Write(transaction *banktransfer.Transaction) error
}

type writer struct {
	writer          *kafka.Writer
	writeContext    context.Context
	writeCancelFunc context.CancelFunc
}

func NewTransactionWriter() *writer {
	writeContext, writeCancelFunc := context.WithCancel(context.Background())
	return &writer{&kafka.Writer{
		Addr:     kafka.TCP(connect),
		Topic:    Topic,
		Balancer: &kafka.LeastBytes{},
	}, writeContext, writeCancelFunc}
}

func (w *writer) Write(transaction *banktransfer.Transaction) error {
	bytes, err := json.Marshal(&transaction)
	if err != nil {
		return fmt.Errorf("error marshalling transaction: %w", err)
	}
	if err := w.writer.WriteMessages(w.writeContext, kafka.Message{
		Key:   []byte(transaction.Id),
		Value: bytes,
	}); err != nil {
		return fmt.Errorf("can't write message to kafka server: %w", err)
	}
	return nil
}

func (w *writer) Close() error {
	w.writeCancelFunc()
	if err := w.writer.Close(); err != nil {
		return fmt.Errorf("couldn't close connection to Kafka server: %w", err)
	}
	return nil
}
