package communication

import (
	"context"
	"crawler/ds"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaTopicWriter[TKey any, TValue any] struct {
	brokersList []string
	topic       string
	logger      *log.Logger
}

func NewKafkaWriter[TKey any, TValue any](brokersList []string, topic string, logger *log.Logger) KafkaTopicWriter[TKey, TValue] {
	return KafkaTopicWriter[TKey, TValue]{
		brokersList: brokersList,
		topic:       topic,
		logger:      logger,
	}
}

func (k KafkaTopicWriter[TKey, TValue]) WriteJSON(ctx context.Context, kvp ds.KeyValuePair[TKey, TValue]) error {
	keyBytes, err := json.Marshal(kvp.Key)
	if err != nil {
		return fmt.Errorf("error while marshaling key: %w", err)
	}

	valueBytes, err := json.Marshal(kvp.Value)
	if err != nil {
		return fmt.Errorf("error while marshaling value: %w", err)
	}

	return k.WriteBytes(ctx, ds.NewKeyValuePair(keyBytes, valueBytes))
}

func (k KafkaTopicWriter[TKey, TValue]) WriteBytes(ctx context.Context, kvp ds.KeyValuePair[[]byte, []byte]) error {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(k.brokersList...),
		Topic:                  k.topic,
		Async:                  true,
		AllowAutoTopicCreation: true,
	}

	defer func() {
		if err := w.Close(); err != nil {
			k.logger.Printf("ERROR: cannot close kafka writer: %s", err)
		}
	}()

	return w.WriteMessages(ctx, kafka.Message{
		Key:   kvp.Key,
		Value: kvp.Value,
	})
}
