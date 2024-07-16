package cronjob

import (
	"context"
	"testing"

	cron "github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	nopLogger := zerolog.Nop()
	cfg := Config{
		Logger: &nopLogger,
	}

	cronjob := New(cfg)
	assert.NotNil(t, cronjob, "Expected cronjob to be not nil")
}

func TestAddJob(t *testing.T) {
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
		{
			name: "Valid cron spec",
			fields: fields{
				instance: cron.New(),
			},
			args: args{
				spec: "* * * * *",
				job:  func() {},
			},
			wantErr: false,
		},
		{
			name: "Invalid cron spec",
			fields: fields{
				instance: cron.New(),
			},
			args: args{
				spec: "invalid spec",
				job:  func() {},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			c := &cronjobImpl{
				instance: tt.fields.instance,
			}
			if err := c.AddJob(context.Background(), tt.args.spec, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("cronjobImpl.AddJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStart(t *testing.T) {
	type fields struct {
		instance *cron.Cron
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Start cron job",
			fields: fields{
				instance: cron.New(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			c := &cronjobImpl{
				instance: tt.fields.instance,
			}
			c.Start()
			// Add assertions to check if the cron job started
		})
	}
}

func TestStop(t *testing.T) {
	type fields struct {
		instance *cron.Cron
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Stop cron job",
			fields: fields{
				instance: cron.New(),
			},
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			c := &cronjobImpl{
				instance: tt.fields.instance,
			}
			c.Stop()
			// Add assertions to check if the cron job stopped
		})
	}
}
