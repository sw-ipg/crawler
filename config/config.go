package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Settings struct {
	ServerPort               int      `json:"server_port" mapstructure:"server_port"`
	KafkaBrokersList         []string `json:"kafka_brokers_list" mapstructure:"kafka_brokers_list"`
	KafkaUrlsTopic           string   `json:"kafka_urls_topic" mapstructure:"kafka_urls_topic"`
	KafkaResultsTopic        string   `json:"kafka_results_topic" mapstructure:"kafka_results_topic"`
	KafkaConsumerGroupId     string   `json:"kafka_consumer_group_id" mapstructure:"kafka_consumer_group_id"`
	PostgresConnectionString string   `json:"postgres_connection_string" mapstructure:"postgres_connection_string"`
}

func ReadConfig() (s Settings, err error) {
	viper.AddConfigPath("./")
	viper.SetConfigType("json")
	viper.SetConfigName("config.json")
	if err = viper.ReadInConfig(); err != nil {
		return s, fmt.Errorf("cannot viper read: %w", err)
	}

	if err = viper.Unmarshal(&s); err != nil {
		return s, fmt.Errorf("cannot unmarshal to struct: %w", err)
	}
	return s, nil
}
