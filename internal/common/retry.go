package common

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

func Retry(run func() error, exitCond func(err error) bool) error {
	for {
		err := run()
		if err == nil {
			return nil
		}

		if exitCond != nil && exitCond(err) {
			return err
		}
	}
}

const (
	RetryMaxWait  = time.Minute
	RetryInterval = time.Second * 5
)

// ported from dockertest package
func ExponentialBackoffRetry(run func() error) error {
	bo := backoff.NewExponentialBackOff()
	bo.MaxInterval = RetryInterval
	bo.MaxElapsedTime = RetryMaxWait

	return backoff.Retry(run, bo) //nolint:wrapcheck
}
