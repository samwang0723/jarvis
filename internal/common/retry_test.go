package common_test

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/samwang0723/jarvis/internal/common"
)

func TestRetry(t *testing.T) {
	t.Parallel()

	type args struct {
		run      func() error
		exitCond func(err error) bool
	}

	var (
		errSome      = errors.New("some error")
		retriedCount int32
	)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "run without error",
			args: args{
				run: func() error {
					return nil
				},
			},
		},
		{
			name: "retry until success",
			args: args{
				run: func() error {
					if atomic.LoadInt32(&retriedCount) > 3 {
						return nil
					}

					atomic.AddInt32(&retriedCount, 1)

					return errSome
				},
			},
		},
		{
			name: "retry until expected errors",
			args: args{
				run: func() error {
					if atomic.LoadInt32(&retriedCount) > 2 {
						return errSome
					}

					atomic.AddInt32(&retriedCount, 1)

					return errors.New("unexpected error")
				},
				exitCond: func(err error) bool {
					return errors.Is(err, errSome)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := common.Retry(tt.args.run, tt.args.exitCond)
			if tt.wantErr && err == nil {
				t.Error("expect error got nil")
			}
		})
	}
}
