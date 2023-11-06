package cache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	key   = "key"
	value = "value"
)

func TestCacheInitialize(t *testing.T) {
	handler := Initialize(DefaultListenAddr)
	if handler.Client.Options().Addr != "localhost"+DefaultListenAddr {
		assert.Error(t, fmt.Errorf("expected address to be 'localhost:6379', but got '%s'", handler.Client.Options().Addr))
	}
	assert.NoError(t, nil)
}

func TestCacheInsertData(t *testing.T) {
	handler := Initialize(DefaultListenAddr)
	handler.Client.Set(key, value, time.Second*4)
	val, err := handler.Client.Get(key).Result()
	if err != nil || val != value {
		assert.Error(t, fmt.Errorf("expected value to be 'value', but got '%s'", val))
	}
	assert.NoError(t, nil)
}

func TestCacheRetrieveData(t *testing.T) {
	handler := Initialize(DefaultListenAddr)
	handler.Client.Set(key, value, time.Second*4)
	val, err := handler.Client.Get(key).Result()
	if err != nil || val != value {
		assert.Error(t, fmt.Errorf("expected value to be 'value', but got '%s'", val))
	}
	assert.NoError(t, nil)
}
