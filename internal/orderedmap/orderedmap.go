package orderedmap

import "sync"

type KeyValuePair struct {
	Key   string
	Value string
}

type OrderedMap struct {
	mu    sync.RWMutex // Assume this mutex is used for thread-safe access to the map
	items map[string]KeyValuePair
	order []string
}

// For initializing the OrderedMap struct when we start a server
func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		items: make(map[string]KeyValuePair),
		order: make([]string, 0),
	}
}

func (om *OrderedMap) Add(key, value string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, exists := om.items[key]; !exists {
		om.order = append(om.order, key)
	}
	om.items[key] = KeyValuePair{Key: key, Value: value}
}

func (om *OrderedMap) Get(key string) (KeyValuePair, bool) {
	om.mu.Lock()
	defer om.mu.Unlock()

	val, exists := om.items[key]
	return val, exists
}

func (om *OrderedMap) Delete(key string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if _, exists := om.items[key]; exists {
		delete(om.items, key)
		for i, k := range om.order {
			if k == key {
				om.order = append(om.order[:i], om.order[i+1:]...)
				break
			}
		}
	}
}

func (om *OrderedMap) GetAll() []KeyValuePair {
	om.mu.RLock() // Lock for reading since we are iterating over shared structure
	defer om.mu.RUnlock()

	var wg sync.WaitGroup
	results := make([]KeyValuePair, len(om.order))

	for i, key := range om.order {
		wg.Add(1)
		// Launch a goroutine for each key
		go func(i int, key string) {
			defer wg.Done()
			// No need to lock here again if om.items read access is inherently thread-safe
			// Otherwise, use om.mu.RLock() and om.mu.RUnlock() to protect each read
			value := om.items[key]
			results[i] = KeyValuePair{Key: key, Value: value.Value}
		}(i, key)
	}

	wg.Wait() // Wait for all goroutines to finish

	return results
}
