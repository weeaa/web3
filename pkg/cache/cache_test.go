package cache

import (
	"testing"
	"time"
)

const (
	key   = "key"
	value = "value"
)

func TestCacheInitialize(t *testing.T) {
	handler := Initialize(DefaultPort)
	if handler.Client.Options().Addr != "localhost"+DefaultPort {
		t.Errorf("expected address to be 'localhost:6379', but got '%s'", handler.Client.Options().Addr)
	}
}

func TestCacheInsertData(t *testing.T) {
	handler := Initialize(DefaultPort)
	handler.Client.Set(key, value, time.Second*4)
	val, err := handler.Client.Get(key).Result()
	if err != nil || val != value {
		t.Errorf("expected value to be 'value', but got '%s'", val)
	}
}

func TestCacheRetrieveData(t *testing.T) {
	handler := Initialize(DefaultPort)
	handler.Client.Set(key, value, time.Second*4)
	val, err := handler.Client.Get(key).Result()
	if err != nil || val != value {
		t.Errorf("expected value to be 'value', but got '%s'", val)
	}
}
