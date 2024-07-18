# jarvis
Taiwan Stock Analysis and Buy Selection Core System.

## Setup project
```bash
./bin/setup.sh
```

## Setup Docker Environment
```bash
docker-compose -f build/docker/postgresql/docker-compose.yml up
docker-compose -f build/docker/redis/docker-compose.yml up
docker-compose -f build/docker/kafka/docker-compose.yml up
```

### Execute SQL Migration
```bash
make db-pg-init-main
make db-pg-migrate
```

## Start Application
```
$ docker-compose -p jarvis -f build/docker/app/docker-compose.yml up
```

### Generate Protobuf
This project is using protobuf with grpc-gatway to have both grpc & http served.
https://github.com/grpc-ecosystem/grpc-gateway

```bash
make proto
```
