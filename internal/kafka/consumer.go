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

	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/helper"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
	"github.com/segmentio/kafka-go"
)

// Config encapsulates the settings for configuring the redis service.
type Config struct {
	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *zerolog.Logger

	GroupID string
	Brokers []string
	Topics  []string
}

type kafkaImpl struct {
	instance Reader
	cfg      Config
}

const (
	queueCapacity    = 1024
	sessionTimeout   = 10 * time.Second
	rebalanceTimeout = 5 * time.Second
	maxWait          = 1 * time.Second
	minBytes         = 1    // 1B
	maxBytes         = 10e6 // 10MB
)

//go:generate mockgen -source=consumer.go -destination=mocks/kafka.go -package=kafka
type Reader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

//nolint:nolintlint, gomnd
func New(cfg Config, reader Reader) ikafka.IKafka {
	if reader == nil {
		reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:          cfg.Brokers,
			GroupTopics:      cfg.Topics,
			GroupID:          cfg.GroupID, // having consumer group id to prevent duplication of message consumption
			QueueCapacity:    queueCapacity,
			SessionTimeout:   sessionTimeout,
			RebalanceTimeout: rebalanceTimeout,
			MaxWait:          maxWait,
			MinBytes:         minBytes,
			MaxBytes:         maxBytes,
			Dialer: &kafka.Dialer{
				Timeout:       10 * time.Second,
				KeepAlive:     30 * time.Second,
				DualStack:     true,
				FallbackDelay: 10 * time.Millisecond,
			},
		})
	}
	return &kafkaImpl{
		instance: reader,
		cfg:      cfg,
	}
}

func (k *kafkaImpl) ReadMessage(ctx context.Context) (ikafka.ReceivedMessage, error) {
	msg, err := k.instance.ReadMessage(ctx)

	k.cfg.Logger.Info().
		Msgf("Kafka:ReadMessage: read data: %s, err: %s", helper.Bytes2String(msg.Value), err)

	return ikafka.ReceivedMessage{
		Topic:   msg.Topic,
		Message: msg.Value,
	}, err
}

func (k *kafkaImpl) Close() error {
	k.cfg.Logger.Info().Msg("Kafka:Close")

	err := k.instance.Close()
	if err != nil {
		k.cfg.Logger.Error().Msgf("Kafka:Close: Close failed: %s", err)
	}

	return err
}
