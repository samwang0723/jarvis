package testhelper

import config "github.com/samwang0723/jarvis/configs"

func LoadTestConfig() {
	testConfig := &config.Config{
		Log: struct {
			Level string `yaml:"level"`
		}{
			Level: "debug",
		},
		JwtSecret:       "test-jwt-secret",
		RecaptchaSecret: "test-recaptcha-secret",
		Kafka: struct {
			GroupID string   `yaml:"groupId"`
			Brokers []string `yaml:"brokers"`
			Topics  []string `yaml:"topics"`
		}{
			GroupID: "test-group",
			Brokers: []string{"localhost:9092"},
			Topics:  []string{"test-topic"},
		},
		RedisCache: struct {
			Master        string   `yaml:"master"`
			Password      string   `yaml:"password"`
			SentinelAddrs []string `yaml:"sentinelAddrs"`
		}{
			Master:   "localhost:6379",
			Password: "test-redis-password",
		},
		Server: struct {
			Name     string `yaml:"name"`
			Host     string `yaml:"host"`
			Version  string `yaml:"version"`
			Port     int    `yaml:"port"`
			GrpcPort int    `yaml:"grpcPort"`
		}{
			Name:     "test-server",
			Host:     "localhost",
			Version:  "1.0.0",
			Port:     8080,
			GrpcPort: 50051,
		},
		Database: struct {
			User         string `yaml:"user"`
			Password     string `yaml:"password"`
			Host         string `yaml:"host"`
			Database     string `yaml:"database"`
			Port         string `yaml:"port"`
			MaxLifetime  int    `yaml:"maxLifetime"`
			MaxIdleConns int    `yaml:"maxIdleConns"`
			MaxOpenConns int    `yaml:"maxOpenConns"`
		}{
			User:         "test-user",
			Password:     "test-password",
			Host:         "localhost",
			Database:     "test-db",
			Port:         "5432",
			MaxLifetime:  3600,
			MaxIdleConns: 10,
			MaxOpenConns: 100,
		},
	}

	// Set the test config as the current config
	*config.GetCurrentConfig() = *testConfig
}
