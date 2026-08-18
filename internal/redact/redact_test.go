package redact

import (
	"reflect"
	"testing"
)

func TestValueRedactsSensitiveKeysAndConfiguredSecrets(t *testing.T) {
	input := map[string]any{
		"password": "visible", "message": "prefix sk-live suffix",
		"nested": []any{map[string]any{"api_key": "also-visible"}},
	}
	want := map[string]any{
		"password": "[REDACTED]", "message": "prefix [REDACTED] suffix",
		"nested": []any{map[string]any{"api_key": "[REDACTED]"}},
	}
	if got := New([]string{"sk-live"}).Value(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("Value() = %#v, want %#v", got, want)
	}
}
