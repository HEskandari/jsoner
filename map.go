package jsoner

import "sync"

type Map struct {
	lock sync.RWMutex
	data map[interface{}]interface{}
}

// NewMap creates a thread safe map
func NewMap() *Map {
	return &Map{
		data: make(map[interface{}]interface{}, 32),
	}
}

// Load is same as sync.Map Load
func (m *Map) Load(key interface{}) (elem interface{}, found bool) {
	m.lock.RLock()
	elem, found = m.data[key]
	m.lock.RUnlock()
	return
}

// Store is same as sync.Map Store
func (m *Map) Store(key interface{}, elem interface{}) {
	m.lock.Lock()
	m.data[key] = elem
	m.lock.Unlock()
}