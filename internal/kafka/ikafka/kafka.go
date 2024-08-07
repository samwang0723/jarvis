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
package ikafka

import (
	"context"
)

const (
	DailyClosesV1        = "dailycloses-v1"
	StocksV1             = "stocks-v1"
	ThreePrimaryV1       = "threeprimary-v1"
	StakeConcentrationV1 = "stakeconcentration-v1"
)

type IKafka interface {
	ReadMessage(ctx context.Context) (ReceivedMessage, error)
	Close() error
}

type ReceivedMessage struct {
	Topic   string
	Message []byte
}
