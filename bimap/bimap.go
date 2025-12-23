package bimap

import (
	"sort"
	"sync"
)

// BiMap is a thread-safe bidirectional map structure
type BiMap struct {
	mu         sync.RWMutex
	keyToValue map[string]string
	valueToKey map[string]string
	keys       []string // Ordered keys for deterministic iteration
}

// NewBiMap creates a new BiMap
func NewBiMap() *BiMap {
	return &BiMap{
		keyToValue: make(map[string]string),
		valueToKey: make(map[string]string),
		keys:       make([]string, 0),
	}
}

// Set adds a key-value pair to the BiMap (thread-safe)
func (bm *BiMap) Set(key string, value string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// Remove old mappings if they exist
	if oldValue, exists := bm.keyToValue[key]; exists {
		delete(bm.valueToKey, oldValue)
	} else {
		// New key, add to ordered list
		bm.keys = append(bm.keys, key)
		sort.Strings(bm.keys) // Keep sorted for deterministic order
	}

	if oldKey, exists := bm.valueToKey[value]; exists {
		delete(bm.keyToValue, oldKey)
		bm.removeKeyFromList(oldKey)
	}

	bm.keyToValue[key] = value
	bm.valueToKey[value] = key
}

// Get returns the value for a given key (thread-safe)
func (bm *BiMap) Get(key string) (string, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	value, exists := bm.keyToValue[key]
	return value, exists
}

// GetByValue returns the key for a given value (thread-safe)
func (bm *BiMap) GetByValue(value string) (string, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	key, exists := bm.valueToKey[value]
	return key, exists
}

// GetValuesWithRange returns key-value pairs from the BiMap with pagination
// Uses ordered keys for deterministic results
func (bm *BiMap) GetValuesWithRange(start, limit int) map[string]string {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	result := make(map[string]string)

	if start >= len(bm.keys) {
		return result
	}

	end := min(start + limit, len(bm.keys))

	for i := start; i < end; i++ {
		key := bm.keys[i]
		result[key] = bm.keyToValue[key]
	}

	return result
}

// Delete removes a key-value pair from the BiMap (thread-safe)
func (bm *BiMap) Delete(key string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if value, exists := bm.keyToValue[key]; exists {
		delete(bm.keyToValue, key)
		delete(bm.valueToKey, value)
		bm.removeKeyFromList(key)
	}
}

// Len returns the number of mappings (thread-safe)
func (bm *BiMap) Len() int {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return len(bm.keyToValue)
}

// removeKeyFromList removes a key from the ordered keys list
// Must be called with lock held
func (bm *BiMap) removeKeyFromList(key string) {
	for i, k := range bm.keys {
		if k == key {
			bm.keys = append(bm.keys[:i], bm.keys[i+1:]...)
			break
		}
	}
}