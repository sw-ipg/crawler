package communication

import (
	"context"
	"crawler/ds"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaTopicReader[TKey any, TValue any] struct {
	brokersList    []string
	topic, groupId string
	logger         *log.Logger
}

func NewKafkaReader[TKey any, TValue any](brokersList []string, topic, groupId string, logger *log.Logger) KafkaTopicReader[TKey, TValue] {
	return KafkaTopicReader[TKey, TValue]{
		brokersList: brokersList,
		topic:       topic,
		groupId:     groupId,
		logger:      logger,
	}
}

func (k KafkaTopicReader[TKey, TValue]) ReadJSON(ctx context.Context) <-chan ds.KeyValuePair[TKey, TValue] {
	outputChan := make(chan ds.KeyValuePair[TKey, TValue], 1)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: k.brokersList,
		Topic:   k.topic,
		GroupID: k.groupId,
	})

	go func() {
		defer func() {
			close(outputChan)
			if err := reader.Close(); err != nil {
				k.logger.Printf("ERROR: unable to close underlying kafka reader: %s", err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				m, err := reader.ReadMessage(ctx)
				if err != nil {
					k.logger.Printf("ERROR: cannot read message from kafka, consumer shutdown. %s", err)
					return
				}

				var kvp ds.KeyValuePair[TKey, TValue]
				if err = json.Unmarshal(m.Key, &kvp.Key); err != nil {
					k.logger.Printf("ERROR: cannot unmarshal key from kafka topic: %s", err)
					continue
				}

				if err = json.Unmarshal(m.Value, &kvp.Value); err != nil {
					k.logger.Printf("ERROR: cannot unmarshal value from kafka topic: %s", err)
					continue
				}

				outputChan <- kvp
			}
		}
	}()

	return outputChan
}
