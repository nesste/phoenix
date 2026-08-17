package corpus

import "testing"

func TestCanonicalOrdersKeysAndRejectsTrailingDocuments(t *testing.T) {
	value, err := DecodeJSON([]byte(`{"z":1,"a":"line\n","😀":true}`))
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := Canonical(value)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(canonical), `{"a":"line\n","z":1,"😀":true}`; got != want {
		t.Fatalf("Canonical() = %s, want %s", got, want)
	}
	if _, err := DecodeJSON([]byte(`{} {}`)); err == nil {
		t.Fatal("DecodeJSON accepted a second document")
	}
	if _, err := DecodeJSON([]byte(`{"same":1,"same":2}`)); err == nil {
		t.Fatal("DecodeJSON accepted a duplicate object key")
	}
}

func TestCanonicalRejectsUnsupportedNumberForms(t *testing.T) {
	for _, input := range []string{`1.5`, `1e2`, `9223372036854775808`} {
		value, err := DecodeJSON([]byte(input))
		if err != nil {
			t.Fatalf("DecodeJSON(%s): %v", input, err)
		}
		if _, err := Canonical(value); err == nil {
			t.Fatalf("Canonical(%s) succeeded", input)
		}
	}
}
