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
cd db/migration
goose mysql "jarvis:password@tcp(localhost:3306)/jarvis?charset=utf8" up
```
