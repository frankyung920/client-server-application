package server

import (
	"client-server-application/internal/orderedmap"
	"log"
	"strings"
	"sync"
)

const (
	ADD_ITEM    = "addItem"
	DELETE_ITEM = "deleteItem"
	GET_ITEM    = "getItem"
	GET_ALL     = "getAllItems"
)

type OrderedMapCommandHandler struct {
	om *orderedmap.OrderedMap
	mu *sync.Mutex
}

func NewOrderedMapCommandHandler(om *orderedmap.OrderedMap, mu *sync.Mutex) *OrderedMapCommandHandler {
	return &OrderedMapCommandHandler{om: om, mu: mu}
}

func (h *OrderedMapCommandHandler) HandleCommand(command string) {
	action, key, value := readCommandToArgs(command)
	h.mu.Lock()
	defer h.mu.Unlock()

	switch action {
	case ADD_ITEM:
		if key == "" {
			log.Println("Key cannot be empty")
			break
		}
		h.om.Add(key, value)
		log.Printf("Added item %s:%s\n", key, value)
	case DELETE_ITEM:
		if key == "" {
			log.Println("Key cannot be empty")
			break
		}
		h.om.Delete(key)
		log.Printf("Deleted item %s\n", key)
	case GET_ITEM:
		if key == "" {
			log.Println("Key cannot be empty")
			break
		}
		value, exists := h.om.Get(key)
		if exists {
			log.Printf("Item: %s = %s\n", key, value.Value)
		}
	case GET_ALL:
		items := h.om.GetAll()
		for _, item := range items {
			log.Printf("Item: %s = %s\n", item.Key, item.Value)
		}
	default:
		log.Println("Unknown command")
	}
}

// readCommandToArgs is a function that takes a command string and splits it into action, key, and value parts.
// It returns the action, key, and value as string.
// If the command contains arguments, they are extracted and assigned to the key and value variables.
// The key and value strings are trimmed of leading and trailing spaces.
// Any occurrences of the character '’' in the key and value strings are replaced with an empty string.
func readCommandToArgs(command string) (string, string, string) {
	parts := strings.Split(command, "(")
	action := parts[0]
	key := ""
	value := ""
	if len(parts) > 1 {
		args := strings.TrimRight(parts[1], ")")
		argParts := strings.Split(args, ",")
		key = strings.ReplaceAll(strings.TrimSpace(argParts[0]), "’", "")
		if len(argParts) > 1 {
			value = strings.ReplaceAll(strings.TrimSpace(argParts[1]), "’", "")
		}
	}

	return action, key, value
}
