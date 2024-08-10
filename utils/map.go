package utils

import (
	"strconv"
	"strings"
)

// PickByNestedKey find a value in nested map
// for nested key need set every key in one string and seperated with dot.
// For example "foo.bar.x.y.z" this is mean you need z key value.
func PickByNestedKey[K comparable, V any](in map[K]V, key string) any {
	if key == "" {
		return nil
	}

	keys := strings.Split(key, ".")

	var traverse func(data any, keyParts []string) any
	traverse = func(data any, keyParts []string) any {
		if len(keyParts) == 0 {
			return data
		}

		currentKey := keyParts[0]
		remainingKeys := keyParts[1:]

		switch d := data.(type) {
		case map[K]V:
			var key K
			if keyStr, ok := any(currentKey).(K); ok {
				key = keyStr
			} else {
				return nil
			}

			if val, ok := d[key]; ok {
				return traverse(val, remainingKeys)
			}
		case map[string]any:
			if val, ok := d[currentKey]; ok {
				return traverse(val, remainingKeys)
			}
		case map[int]any:
			idx, err := strconv.Atoi(currentKey)
			if err != nil {
				return nil
			}
			if val, ok := d[idx]; ok {
				return traverse(val, remainingKeys)
			}
		case []map[string]interface{}:
			idx, err := strconv.Atoi(currentKey)
			if err != nil || idx < 0 || idx >= len(d) {
				return nil
			}
			return traverse(d[idx], remainingKeys)
		case []any:
			idx, err := strconv.Atoi(currentKey)
			if err != nil || idx < 0 || idx >= len(d) {
				return nil
			}
			return traverse(d[idx], remainingKeys)
		}
		return nil
	}

	return traverse(in, keys)
}
