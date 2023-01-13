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
package cache

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	redismock "github.com/go-redis/redismock/v8"
	"github.com/rs/zerolog/log"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestSetExpire(t *testing.T) {
	t.Parallel()

	type args struct {
		key     string
		expired time.Time
		err     error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Redis SetExpire successfully",
			args: args{
				key:     "test",
				expired: time.Now(),
				err:     nil,
			},
			wantErr: false,
		},
	}

	logger := log.With().Str("test", "redis").Logger()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()

			client, mock := redismock.NewClientMock()

			impl := &redisImpl{
				instance: client,
				cfg: Config{
					Logger: &logger,
				},
			}

			mock.Regexp().ExpectExpireAt(tt.args.key, tt.args.expired).SetErr(tt.args.err)
			impl.SetExpire(ctx, tt.args.key, tt.args.expired)
		})
	}
}

func TestSAdd(t *testing.T) {
	t.Parallel()

	type args struct {
		key   string
		value []string
		err   error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Redis SAdd successfully",
			args: args{
				key:   "test",
				value: []string{"test"},
				err:   nil,
			},
			wantErr: false,
		},
	}

	logger := log.With().Str("test", "redis").Logger()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()

			client, mock := redismock.NewClientMock()

			impl := &redisImpl{
				instance: client,
				cfg: Config{
					Logger: &logger,
				},
			}

			mock.Regexp().ExpectSAdd(tt.args.key, tt.args.value).SetErr(tt.args.err)
			impl.SAdd(ctx, tt.args.key, tt.args.value)
		})
	}
}

func TestSMembers(t *testing.T) {
	t.Parallel()

	type args struct {
		key string
		err error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Redis SMembers successfully",
			args: args{
				key: "test",
				err: nil,
			},
			wantErr: false,
		},
	}

	logger := log.With().Str("test", "redis").Logger()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.TODO()

			client, mock := redismock.NewClientMock()

			impl := &redisImpl{
				instance: client,
				cfg: Config{
					Logger: &logger,
				},
			}

			mock.Regexp().ExpectSMembers(tt.args.key).SetVal([]string{"test"})
			res, err := impl.SMembers(ctx, tt.args.key)

			if (err != nil || len(res) == 0) != tt.wantErr {
				t.Errorf("SMembers() error = %v, res = %v", err, res)
			}
		})
	}
}
