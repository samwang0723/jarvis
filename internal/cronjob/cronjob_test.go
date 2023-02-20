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

package cronjob

import (
	"context"
	"flag"
	"os"
	"testing"

	cron "github.com/robfig/cron/v3"
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

func Test_cronjobImpl_AddJob(t *testing.T) {
	t.Parallel()

	type fields struct {
		instance *cron.Cron
	}
	type args struct {
		spec string
		job  func()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := &cronjobImpl{
				instance: tt.fields.instance,
			}
			if err := c.AddJob(context.Background(), tt.args.spec, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("cronjobImpl.AddJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
