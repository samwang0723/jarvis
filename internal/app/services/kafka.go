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
package services

import (
	"context"

	"github.com/samwang0723/jarvis/internal/app/entity"
	"github.com/samwang0723/jarvis/internal/kafka/ikafka"
	log "github.com/samwang0723/jarvis/internal/logger"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type data struct {
	values *[]interface{}
	topic  string
}

func (s *serviceImpl) ListeningKafkaInput(ctx context.Context) {
	respChan := make(chan data)
	go func() {
		for {
			msg, err := s.consumer.ReadMessage(ctx)
			if err != nil {
				log.Errorf("Kafka:ReadMessage error: %w", err)
				return
			}

			entity, err := unmarshalMessageToEntity(msg)
			if err != nil {
				log.Errorf("Unmarshal (%s) failed: %w", msg.Topic, err)
				return
			}
			respChan <- data{
				topic:  msg.Topic,
				values: &[]interface{}{entity},
			}

			select {
			case <-ctx.Done():
				log.Warn("ListeningKafkaInput: context cancel")
				return
			default:
			}
		}
	}()

	// handler goroutine to insert message from Kafka to database
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Warn("ListeningKafkaInput(respChan): context cancel")
				return
			case obj, ok := <-respChan:
				if ok {
					switch obj.topic {
					case ikafka.DailyClosesV1:
						s.BatchUpsertDailyClose(ctx, obj.values)
					case ikafka.StakeConcentrationV1:
						s.BatchUpsertStakeConcentration(ctx, obj.values)
					case ikafka.StocksV1:
						s.BatchUpsertStocks(ctx, obj.values)
					case ikafka.ThreePrimaryV1:
						s.BatchUpsertThreePrimary(ctx, obj.values)
					}
				}
			}
		}
	}()

}

func (s *serviceImpl) StopKafka() error {
	return s.consumer.Close()
}

func unmarshalMessageToEntity(msg ikafka.ReceivedMessage) (interface{}, error) {
	var err error
	var output interface{}
	switch msg.Topic {
	case ikafka.DailyClosesV1:
		var obj entity.DailyClose
		err = json.Unmarshal(msg.Message, &obj)
		output = &obj
	case ikafka.StakeConcentrationV1:
		var obj entity.StakeConcentration
		err = json.Unmarshal(msg.Message, &obj)
		output = &obj
	case ikafka.StocksV1:
		var obj entity.Stock
		err = json.Unmarshal(msg.Message, &obj)
		output = &obj
	case ikafka.ThreePrimaryV1:
		var obj entity.ThreePrimary
		err = json.Unmarshal(msg.Message, &obj)
		output = &obj
	}
	return output, err
}
