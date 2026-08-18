// Package jsonptr resolves RFC 6901 pointers over normalized JSON values.
package jsonptr

import (
	"strconv"
	"strings"
)

func Valid(pointer string) bool {
	if pointer == "" {
		return true
	}
	if !strings.HasPrefix(pointer, "/") {
		return false
	}
	for _, token := range strings.Split(strings.TrimPrefix(pointer, "/"), "/") {
		for index := 0; index < len(token); index++ {
			if token[index] == '~' && (index+1 >= len(token) || (token[index+1] != '0' && token[index+1] != '1')) {
				return false
			}
			if token[index] == '~' {
				index++
			}
		}
	}
	return true
}

func Tokens(pointer string) []string {
	if pointer == "" {
		return nil
	}
	encoded := strings.Split(strings.TrimPrefix(pointer, "/"), "/")
	for index := range encoded {
		encoded[index] = strings.ReplaceAll(strings.ReplaceAll(encoded[index], "~1", "/"), "~0", "~")
	}
	return encoded
}

func Resolve(value any, pointer string) (any, bool) {
	if !Valid(pointer) {
		return nil, false
	}
	current := value
	for _, token := range Tokens(pointer) {
		switch typed := current.(type) {
		case map[string]any:
			var exists bool
			current, exists = typed[token]
			if !exists {
				return nil, false
			}
		case []any:
			index, err := strconv.Atoi(token)
			if err != nil || index < 0 || index >= len(typed) {
				return nil, false
			}
			current = typed[index]
		default:
			return nil, false
		}
	}
	return current, true
}
