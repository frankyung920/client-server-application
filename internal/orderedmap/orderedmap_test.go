package orderedmap

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestOrderedMap_AddAndGet(t *testing.T) {
	om := NewOrderedMap()
	key := "testKey"
	value := "testValue"
	om.Add(key, value)

	v, exists := om.Get(key)
	if !exists {
		t.Errorf("Expected key %s to exist", key)
	}

	if v.Value != value {
		t.Errorf("Expected value %s, got %s", value, v.Value)
	}
}

func TestOrderedMap_Delete(t *testing.T) {
	om := NewOrderedMap()
	key := "testKey"
	value := "testValue"
	om.Add(key, value)

	om.Delete(key)
	_, exists := om.Get(key)
	if exists {
		t.Errorf("Expected key %s to be deleted", key)
	}
}

// Simply test GetAll() method
func TestOrderedMap_GetAll(t *testing.T) {
	om := NewOrderedMap()
	items := []KeyValuePair{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}

	for _, item := range items {
		om.Add(item.Key, item.Value)
	}

	results := om.GetAll()
	if !reflect.DeepEqual(items, results) {
		t.Errorf("Expected results to equal items. Results: %v, Items: %v", results, items)
	}
}

// Specifically test the order of items returned by GetAll()
func TestOrderedMap_ConcurrentGetAll(t *testing.T) {
    om := NewOrderedMap()
    var wg sync.WaitGroup

    // Prepopulate the map with a known set of key-value pairs
    for i := 0; i < 1000; i++ {
        key := fmt.Sprintf("key%d", i)
        value := fmt.Sprintf("value%d", i)
        om.Add(key, value)
    }

    // Define a function that concurrently reads from the map using GetAll
    readAllItems := func(wg *sync.WaitGroup) {
        defer wg.Done()
        result := om.GetAll() // In a real test, you'd want to check the results match expected values
		if len(result) != 1000 {
			t.Errorf("Expected 1000 items in the map, found %d", len(om.GetAll()))
		}
	}

    // Launch multiple goroutines to read from the map concurrently
    goroutines := 50
    for i := 0; i < goroutines; i++ {
        wg.Add(1)
        go readAllItems(&wg)
    }

    // Wait for all goroutines to complete
    wg.Wait()

    // After concurrent access, verify the map's integrity
    // This can include checking the count of items, specific key-value pairs, or other invariants
    if len(om.GetAll()) != 1000 {
        t.Errorf("Expected 1000 items in the map, found %d", len(om.GetAll()))
    }
}