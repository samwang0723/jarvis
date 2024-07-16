package remotetest

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaContainer struct {
	Address string
	Topics  sync.Map
}

func CreateKafkaContainer() (*KafkaContainer, error) {
	address := "localhost:19092"

	if envAddress, ok := os.LookupEnv("TEST_KAFKA_ADDRESS"); ok {
		address = envAddress
	}

	return &KafkaContainer{
		Address: address,
		Topics:  sync.Map{},
	}, nil
}

func (k *KafkaContainer) Purge() error {
	conn, err := kafka.Dial("tcp", k.Address)
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}

	var controllerConn *kafka.Conn

	controllerConn, err = kafka.Dial(
		"tcp",
		net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)),
	)
	if err != nil {
		return fmt.Errorf("failed to dial kafka controller: %w", err)
	}
	defer controllerConn.Close()

	k.Topics.Range(func(key, _ interface{}) bool {
		topic := key.(string)              // nolint // never fails
		controllerConn.DeleteTopics(topic) // nolint

		return true
	})

	return nil
}

func (k *KafkaContainer) CreateTopic(topic string) error {
	conn, err := kafka.Dial("tcp", k.Address)
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}

	var controllerConn *kafka.Conn

	controllerConn, err = kafka.Dial(
		"tcp",
		net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)),
	)
	if err != nil {
		return fmt.Errorf("failed to dial kafka controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	err = controllerConn.CreateTopics(topicConfig)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	k.Topics.Store(topic, true)

	return nil
}
