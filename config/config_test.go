package config_test

import (
	"savemanager/config"
	"testing"
)

func TestConfigGet(t *testing.T) {
	config.GetConfig()
}

func BenchmarkConfigGet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config.GetConfig()

	}
}

func BenchmarkConfigLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config.GetConfig().Log("Log test number ", i)

	}
}
