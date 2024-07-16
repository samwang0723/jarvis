package remotetest

import (
	"fmt"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

type RedisServer struct {
	mrs *miniredis.Miniredis
}

func SetupRedisServer() (*RedisServer, error) {
	mrs, err := miniredis.Run()
	if err != nil {
		return nil, fmt.Errorf("setup redis server error: %w", err)
	}

	return &RedisServer{
		mrs: mrs,
	}, nil
}

func (rs *RedisServer) SetupRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: rs.mrs.Addr(),
	})

	return client
}

func (rs *RedisServer) Purge() {
	rs.mrs.Close()
}

func (rs *RedisServer) FastForward(d time.Duration) {
	rs.mrs.FastForward(d)
}

func SetupRedisClient(t *testing.T) *redis.Client {
	t.Helper()

	server, err := SetupRedisServer()
	if err != nil {
		t.Fatalf("create redis server error: %s", err)
	}

	t.Cleanup(server.Purge)

	client := server.SetupRedisClient()

	t.Cleanup(func() {
		client.Close()
	})

	return client
}
