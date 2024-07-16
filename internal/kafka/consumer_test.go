package kafka_test

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	k "github.com/samwang0723/jarvis/internal/kafka"
	kafka_mock "github.com/samwang0723/jarvis/internal/kafka/mocks"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zerolog.New(zerolog.Nop())
	cfg := k.Config{
		Logger:  &logger,
		GroupID: "test-group",
		Brokers: []string{"localhost:9092"},
		Topics:  []string{"test-topic"},
	}

	mockReader := kafka_mock.NewMockReader(ctrl)
	k := k.New(cfg, mockReader)
	assert.NotNil(t, k)
}

func TestReadMessage(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := kafka_mock.NewMockReader(ctrl)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	cfg := k.Config{
		Logger:  &logger,
		GroupID: "test-group",
		Brokers: []string{"localhost:9092"},
		Topics:  []string{"test-topic"},
	}

	k := k.New(cfg, mockReader)

	msg := kafka.Message{
		Topic: "test-topic",
		Value: []byte("test message"),
	}

	mockReader.EXPECT().ReadMessage(gomock.Any()).Return(msg, nil)

	receivedMsg, err := k.ReadMessage(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "test-topic", receivedMsg.Topic)
	assert.Equal(t, []byte("test message"), receivedMsg.Message)
}

func TestClose(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := kafka_mock.NewMockReader(ctrl)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	cfg := k.Config{
		Logger:  &logger,
		GroupID: "test-group",
		Brokers: []string{"localhost:9092"},
		Topics:  []string{"test-topic"},
	}

	k := k.New(cfg, mockReader)

	mockReader.EXPECT().Close().Return(nil)

	err := k.Close()
	assert.NoError(t, err)
}
