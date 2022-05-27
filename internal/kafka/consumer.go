// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package kafka

import (
	"context"

	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
	log "github.com/samwang0723/jarvis/internal/logger"
	"github.com/segmentio/kafka-go"
)

type kafkaImpl struct {
	instance *kafka.Reader
}

func New(cfg *config.Config) ikafka.IKafka {
	return &kafkaImpl{
		instance: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     cfg.Kafka.Brokers,
			GroupTopics: cfg.Kafka.Topics,
			GroupID:     cfg.Kafka.GroupId, // having consumer group id to prevent duplication of message consumption
			Partition:   cfg.Kafka.Partition,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
		}),
	}
}

func (k *kafkaImpl) ReadMessage(ctx context.Context) (ikafka.ReceivedMessage, error) {
	msg, err := k.instance.ReadMessage(ctx)
	log.Infof("Kafka:ReadMessage: read data: %+v, err: %s", msg, err)

	return ikafka.ReceivedMessage{
		Topic:   msg.Topic,
		Message: msg.Value,
	}, err
}

func (k *kafkaImpl) Close() error {
	log.Info("Kafka:Close")
	err := k.instance.Close()
	if err != nil {
		log.Errorf("Close failed: %w", err)
	}
	return err
}
