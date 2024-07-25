#/bin/bash

/opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --create --topic stakeconcentration-v1 --replication-factor 1 --partitions 1
echo "topic stakeconcentration-v1 was created"

/opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --create --topic dailycloses-v1 --replication-factor 1 --partitions 1
echo "topic dailycloses-v1 was created"

/opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --create --topic stocks-v1 --replication-factor 1 --partitions 1
echo "topic stocks-v1 was created"

/opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --create --topic threeprimary-v1 --replication-factor 1 --partitions 1
echo "topic threeprimary-v1 was created"

/opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --create --topic download-v1 --replication-factor 1 --partitions 1
echo "topic download-v1 was created"
