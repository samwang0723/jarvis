package remotetest

import (
	"testing"
)

func TestCreateKafka(t *testing.T) {
	t.Parallel()

	kc, err := CreateKafkaContainer()
	if err != nil {
		t.Fatalf("create kafka error: %s", err)
	}

	if err = kc.CreateTopic("test_topic"); err != nil {
		t.Fatalf("create topic error: %s", err)
	}

	if err = kc.Purge(); err != nil {
		t.Errorf("purge kafka error: %s", err)
	}
}
