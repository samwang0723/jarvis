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

package elastic

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	config "github.com/samwang0723/jarvis/configs"
	clog "github.com/samwang0723/jarvis/internal/logger"

	"github.com/olivere/elastic/v7"
)

func Setup(ctx context.Context, cfg *config.Config) *elastic.Client {
	var client *elastic.Client

	elasticUrl := fmt.Sprintf("http://%s:%d", cfg.ElasticSearch.Host, cfg.ElasticSearch.Port)
	clog.Info("Connecting to elastic search on", elasticUrl)

	//create an elastic search client. connect to the running elastic search db
	client, connError := elastic.NewClient(
		elastic.SetURL(elasticUrl),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(time.Duration(cfg.ElasticSearch.HealthCheckInterval)*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if connError != nil {
		panic(connError)
	}

	info, code, err := client.Ping(elasticUrl).Do(ctx)

	if err != nil {
		panic(err)
	}
	clog.Infof("Elasticsearch returned with code %d and version %s", code, info.Version.Number)

	return client
}
