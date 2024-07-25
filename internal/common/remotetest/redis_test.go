package remotetest

import (
	"context"
	"testing"
	"time"
)

func TestRedisContainer(t *testing.T) {
	t.Parallel()

	client := SetupRedisClient(t)

	ctx := context.Background()

	key, value := "key", "value"

	if err := client.SetNX(ctx, key, value, time.Second).Err(); err != nil {
		t.Errorf("set value to redis error: %v", err)
	}

	got, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("read redis error: %s %v", "key", err)
	}

	if value != got {
		t.Errorf("expect %v got %v", value, got)
	}

	duration, err := client.TTL(ctx, key).Result()
	if err != nil {
		t.Errorf("read ttl error: %v", err)
	}

	t.Logf("ttl %v", duration.String())
}
