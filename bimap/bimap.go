package bimap

import (
	"sort"
	"sync"
)

type ItemWithExtra struct {
	Key   string
	Value string
	Extra any
}

type entry struct {
	Value string
	Extra any // optional: int, string, map, set, struct, etc.
}

// BiMap is a thread-safe bidirectional map structure
type BiMap struct {
	mu         sync.RWMutex
    keyToEntry map[string]entry
	valueToKey map[string]string
	keys       []string // Ordered keys for deterministic iteration
}

// NewBiMap creates a new BiMap
func NewBiMap() *BiMap {
	return &BiMap{
		keyToEntry: make(map[string]entry),
		valueToKey: make(map[string]string),
		keys:       make([]string, 0),
	}
}

// Set adds a key-value pair to the BiMap (thread-safe)
func (bm *BiMap) Set(key string, value string) {
	bm.SetWithExtra(key, value, nil)
}

func (bm *BiMap) SetWithExtra(key, value string, extra any) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if oldentry, exists := bm.keyToEntry[key]; exists {
		delete(bm.valueToKey, oldentry.Value)
	} else {
		bm.keys = append(bm.keys, key)
		sort.Strings(bm.keys)
	}

	if oldKey, exists := bm.valueToKey[value]; exists {
		delete(bm.keyToEntry, oldKey)
		bm.removeKeyFromList(oldKey)
	}

	bm.keyToEntry[key] = entry{
		Value: value,
		Extra: extra, // can be nil or any type
	}
	bm.valueToKey[value] = key
}

// Get returns the value for a given key (thread-safe)
func (bm *BiMap) Get(key string) (string, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	entry, exists := bm.keyToEntry[key]
	return entry.Value, exists
}

func (bm *BiMap) GetExtra(key string) (string, any, bool) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	entry, exists := bm.keyToEntry[key]
	if !exists {
		return "", nil, false
	}
	return entry.Value, entry.Extra, true
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
func (bm *BiMap) GetWithRange(start, limit int) []ItemWithExtra {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if start >= len(bm.keys) || limit <= 0 {
		return nil
	}

	end := min(start+limit, len(bm.keys))
	result := make([]ItemWithExtra, 0, end-start)

	for i := start; i < end; i++ {
		key := bm.keys[i]
		value := bm.keyToEntry[key]

		result = append(result, ItemWithExtra{
			Key:   key,
			Value: value.Value,
			Extra: value.Extra, // can be nil
		})
	}

	return result
}

// Delete removes a key-value pair from the BiMap (thread-safe)
func (bm *BiMap) Delete(key string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if entry, exists := bm.keyToEntry[key]; exists {
		delete(bm.keyToEntry, key)
		delete(bm.valueToKey, entry.Value)
		bm.removeKeyFromList(key)
	}
}

// Len returns the number of mappings (thread-safe)
func (bm *BiMap) Len() int {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return len(bm.keyToEntry)
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