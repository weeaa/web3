package safemap

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	key   = "testKey"
	value = "testValue"
)

func TestSafeMap(t *testing.T) {
	sm := New[string, string]()
	sm.Set(key, value)

	result, found := sm.Get(key)
	if !found {
		assert.Error(t, errors.New("expected to find the key, but it wasn't found"))
	}
	if result != value {
		assert.Error(t, fmt.Errorf("expected %s, got %s", value, result))
	}

	sm.Delete(key)
	_, found = sm.Get(key)
	if found {
		assert.Error(t, errors.New("expected the key to be deleted, but it's still found"))
	}

	if sm.Len() != 0 {
		assert.Error(t, fmt.Errorf("expected length to be 0, got %d", sm.Len()))
	}

	sm.Set("key1", "value1")
	sm.Set("key2", "value2")

	var counter int
	sm.ForEach(func(k string, v string) {
		counter++
	})
	if counter != 2 {
		assert.Error(t, fmt.Errorf("expected count to be 2, got %d", counter))
	}

	assert.NoError(t, nil)
}

func BenchmarkSafeMapSet(b *testing.B) {
	sm := New[string, string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Set(key, value)
	}
}

func BenchmarkSafeMapGet(b *testing.B) {
	sm := New[string, string]()
	sm.Set(key, value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = sm.Get(key)
	}
}

func BenchmarkSafeMapDelete(b *testing.B) {
	sm := New[string, string]()
	sm.Set(key, value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sm.Delete(key)
	}
}
