// Package redact removes configured secrets and values under sensitive keys
// from normalized JSON data.
package redact

import (
	"sort"
	"strings"
)

type Redactor struct {
	secrets []string
}

func New(secrets []string) Redactor {
	filtered := make([]string, 0, len(secrets))
	for _, secret := range secrets {
		if secret != "" {
			filtered = append(filtered, secret)
		}
	}
	sort.Slice(filtered, func(i, j int) bool { return len(filtered[i]) > len(filtered[j]) })
	return Redactor{secrets: filtered}
}

func (redactor Redactor) Value(value any) any {
	return redactor.value(value, "")
}

func (redactor Redactor) value(value any, key string) any {
	if key != "" && sensitiveKey(key) && value != nil {
		return "[REDACTED]"
	}
	switch typed := value.(type) {
	case string:
		for _, secret := range redactor.secrets {
			typed = strings.ReplaceAll(typed, secret, "[REDACTED]")
		}
		return typed
	case map[string]any:
		result := make(map[string]any, len(typed))
		for childKey, child := range typed {
			result[childKey] = redactor.value(child, childKey)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for index, child := range typed {
			result[index] = redactor.value(child, key)
		}
		return result
	default:
		return value
	}
}

func sensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(key, "-", "_"))
	for _, fragment := range []string{"password", "secret", "token", "credential", "authorization", "api_key"} {
		if strings.Contains(normalized, fragment) {
			return true
		}
	}
	return false
}
