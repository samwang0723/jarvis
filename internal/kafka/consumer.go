// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package kafka

import (
	"context"
	"time"

	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/helper"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
	log "github.com/samwang0723/jarvis/internal/logger"
	"github.com/segmentio/kafka-go"
)

type kafkaImpl struct {
	instance *kafka.Reader
}

const (
	queueCapacity    = 1024
	sessionTimeout   = 10 * time.Second
	rebalanceTimeout = 5 * time.Second
	maxWait          = 1 * time.Second
	minBytes         = 1    // 1B
	maxBytes         = 10e6 // 10MB
)

func New(cfg *config.Config) ikafka.IKafka {
	return &kafkaImpl{
		instance: kafka.NewReader(kafka.ReaderConfig{
			Brokers:          cfg.Kafka.Brokers,
			GroupTopics:      cfg.Kafka.Topics,
			GroupID:          cfg.Kafka.GroupID, // having consumer group id to prevent duplication of message consumption
			QueueCapacity:    queueCapacity,
			SessionTimeout:   sessionTimeout,
			RebalanceTimeout: rebalanceTimeout,
			MaxWait:          maxWait,
			MinBytes:         minBytes,
			MaxBytes:         maxBytes,
		}),
	}
}

func (k *kafkaImpl) ReadMessage(ctx context.Context) (ikafka.ReceivedMessage, error) {
	msg, err := k.instance.ReadMessage(ctx)
	log.Infof("Kafka:ReadMessage: read data: %s, err: %s", helper.Bytes2String(msg.Value), err)

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
