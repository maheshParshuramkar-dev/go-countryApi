package cache

import "sync"

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

type InMemCache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// creating in mem cache and returning it
func NewCache() *InMemCache {
	return &InMemCache{
		data: make(map[string]interface{}),
	}
}

// to get data from in-mem cache
func (m *InMemCache) Get(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	data, ok := m.data[key]
	return data, ok
}

// to set data into in-mem cache
func (m *InMemCache) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
}
