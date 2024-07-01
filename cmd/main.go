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

package main

import (
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	zl "github.com/rs/zerolog/log"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/app/server"
	"github.com/samwang0723/jarvis/internal/helper"
)

const (
	appName              = "jarvis"
	sentryFlushThreshold = 2 * time.Second
)

func main() {
	config.Load()
	cfg := config.GetCurrentConfig()

	logger := zl.With().Str("app", appName).Logger()

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.Sentry.DSN,
		Environment: cfg.Environment,
		Release:     "v2.0.0",
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug: cfg.Sentry.Debug,
	})
	if err != nil {
		logger.Error().Msgf("sentry.Init: %s", err)
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(sentryFlushThreshold)

	// manually set time zone, docker image may not have preset timezone
	time.Local, err = time.LoadLocation(helper.TimeZone)
	if err != nil {
		logger.Error().Msgf("error loading location '%s': %v\n", helper.TimeZone, err)
	}

	server.Serve(cfg, &logger)
}
