# jarvis
Machine Learning Stock Analysis and Buy Selection

## Setup Docker MySQL

### Start Docker Container

```
$ docker-compose -p mysql -f build/docker/mysql/docker-compose.yml up
```

### Configure Database(master)
```
$ docker exec -it mysql-master mysql -u root -p

CREATE USER 'jarvis'@'%' IDENTIFIED BY 'password';
SELECT host, user FROM mysql.user;

CREATE DATABASE jarvis CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL PRIVILEGES ON jarvis.* TO 'jarvis'@'%';
FLUSH PRIVILEGES;

GRANT REPLICATION SLAVE ON *.* TO ‘jarvis’@‘%’;
```

### Configure Database(slave)
```
$ docker exec -it mysql-slave mysql -u root -p

CHANGE MASTER TO
MASTER_HOST='mysql-master',
MASTER_PORT=3306,
MASTER_USER='jarvis',
MASTER_PASSWORD='password',
MASTER_AUTO_POSITION=1;

START SLAVE;
SHOW SLAVE STATUS \G
```

### Execute SQL Migration

```
$ make migrate
```

#### mysqldump
```
$ docker exec mysql-master /usr/bin/mysqldump -u root --all-databases --triggers --routines --events --set-gtid-purged=OFF > backup.sql
$ docker exec -i mysql-master mysql -u root < backup.sql
```

### Start Application Container
```
$ docker-compose -p jarvis -f build/docker/app/docker-compose.yml up
```

### Generate Protobuf

https://github.com/grpc-ecosystem/grpc-gateway

Cloning google/api annotation files
```
$ mkdir third_party/google/api
$ curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > pb/google/api/annotations.proto
$ curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > pb/google/api/http.proto
$ curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto > pb/google/api/field_behavior.proto
$ curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/httpbody.proto > pb/google/api/httpbody.proto
```

Preparation of generate proto
```
$ go mod tidy
$ go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
$ make proto
```

### Use swagger-ui

Clone swagger-ui static files into `third_party/swagger-ui/`
https://github.com/swagger-api/swagger-ui/tree/master/dist
``` 
$ go get -u github.com/jteeuwen/go-bindata/...
$ go get -u github.com/elazarl/go-bindata-assetfs/...
$ go-bindata --nocompress -pkg swagger -o api/swagger/datafile.go third_party/swagger-ui/...
```

Import "protoc-gen-swagger/options/openapiv2.proto" was not found or had errors.
```
$ cp -R $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/protoc-gen-swagger ./third_party
```
