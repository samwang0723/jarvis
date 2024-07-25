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
package services

import (
	"context"
	"io"

	"github.com/gofrs/uuid/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"github.com/samwang0723/jarvis/internal/app/domain"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
	"golang.org/x/xerrors"
)

//nolint:nolintlint, gochecknoglobals
var json = jsoniter.ConfigCompatibleWithStandardLibrary

type data struct {
	values *[]any
	topic  string
}

// Config encapsulates the settings for configuring the redis service.
type KafkaConfig struct {
	Logger *zerolog.Logger

	GroupID string
	Brokers []string
	Topics  []string
}

func (cfg *KafkaConfig) validate() error {
	if cfg.GroupID == "" {
		return xerrors.Errorf("invalid kafka groupId")
	}

	if len(cfg.Brokers) == 0 {
		return xerrors.Errorf("invalid kafka brokers")
	}

	if len(cfg.Topics) == 0 {
		return xerrors.Errorf("invalid kafka topics")
	}

	return nil
}

//nolint:nolintlint, cyclop
func (s *serviceImpl) ListeningKafkaInput(ctx context.Context) {
	respChan := make(chan data)

	go func() {
		s.logger.Info().Str("component", "kafka").Msg("goroutine starting")
		defer s.logger.Info().Str("component", "kafka").Msg("goroutine exited")
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := s.consumer.ReadMessage(ctx)
				if err != nil {
					if err == context.Canceled || err == io.EOF {
						return
					}
					s.logger.Error().
						Str("component", "kafka").Err(err).
						Msgf("read message")
					continue
				}

				ent, err := unmarshalMessageTodomain(msg)
				if err != nil {
					s.logger.Error().
						Str("component", "kafka").Err(err).
						Msg("domain model unmarshal")
					continue
				}
				respChan <- data{
					topic:  msg.Topic,
					values: &[]any{ent},
				}
			}
		}
	}()

	// handler goroutine to insert message from Kafka to database
	go func() {
		s.logger.Info().Str("component", "handler").Msg("goroutine starting")
		defer s.logger.Info().Str("component", "handler").Msg("goroutine exited")
		for {
			select {
			case <-ctx.Done():
				return
			case obj, ok := <-respChan:
				var err error
				if ok {
					switch obj.topic {
					case ikafka.DailyClosesV1:
						err = s.BatchUpsertDailyClose(ctx, obj.values)
					case ikafka.StakeConcentrationV1:
						err = s.BatchUpsertStakeConcentration(ctx, obj.values)
					case ikafka.StocksV1:
						err = s.BatchUpsertStocks(ctx, obj.values)
					case ikafka.ThreePrimaryV1:
						err = s.BatchUpsertThreePrimary(ctx, obj.values)
					}

					if err != nil {
						s.logger.Error().
							Str("component", "kafka").
							Msgf("batch upsert (%s) failed: %s", obj.topic, err.Error())
					}
				}
			}
		}
	}()
}

func (s *serviceImpl) StopKafka() error {
	return s.consumer.Close()
}

func unmarshalMessageTodomain(msg ikafka.ReceivedMessage) (any, error) {
	var err error

	var output any

	switch msg.Topic {
	case ikafka.DailyClosesV1:
		var obj domain.DailyClose
		err = json.Unmarshal(msg.Message, &obj)
		obj.ID.ID = uuid.Must(uuid.NewV4())
		output = &obj
	case ikafka.StakeConcentrationV1:
		var obj domain.StakeConcentration
		err = json.Unmarshal(msg.Message, &obj)
		obj.ID.ID = uuid.Must(uuid.NewV4())
		output = &obj
	case ikafka.StocksV1:
		var obj domain.Stock
		err = json.Unmarshal(msg.Message, &obj)
		output = &obj
	case ikafka.ThreePrimaryV1:
		var obj domain.ThreePrimary
		err = json.Unmarshal(msg.Message, &obj)
		obj.ID.ID = uuid.Must(uuid.NewV4())
		output = &obj
	}

	return output, err
}
