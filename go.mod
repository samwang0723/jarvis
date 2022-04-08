module github.com/samwang0723/jarvis

go 1.16

require (
	github.com/getsentry/sentry-go v0.12.0
	github.com/go-errors/errors v1.0.1
	github.com/robfig/cron/v3 v3.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/sony/sonyflake v1.0.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f
	golang.org/x/text v0.3.7
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.2
	gorm.io/gorm v1.21.15
)

require (
	github.com/bsm/redislock v0.7.2
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.8.0
	github.com/heptiolabs/healthcheck v0.0.0-20211123025425-613501dd5deb
	github.com/johnbellone/grpc-middleware-sentry v0.2.0
	github.com/joho/godotenv v1.4.0
	github.com/prometheus/client_golang v1.12.1 // indirect
	golang.org/x/sys v0.0.0-20220330033206-e17cdc41300f // indirect
	google.golang.org/genproto v0.0.0-20220304144024-325a89244dc8
	google.golang.org/grpc v1.44.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
	gorm.io/plugin/dbresolver v1.1.0
)
