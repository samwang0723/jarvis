# jarvis
Machine Learning Stock Analysis and Buy Selection

## Setup Docker MySQL

### Start Docker Container

```
docker-compose up
docker ps
```

### Configure Database and Access
```
docker exec -it mysql-master bin/bash
mysql -h 127.0.0.1 -P 3306 -u root

CREATE USER 'jarvis'@'%' IDENTIFIED BY 'password';
SELECT host, user FROM mysql.user;

CREATE DATABASE jarvis CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL PRIVILEGES ON jarvis.* TO 'jarvis'@'%';
FLUSH PRIVILEGES;
```

### Execute SQL Migration

```
$ cd internal/db/migration
$ goose mysql "jarvis:password@tcp(localhost:3306)/jarvis?charset=utf8" up
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
