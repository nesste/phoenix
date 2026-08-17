package corpus

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode/utf16"
)

// DigestPrefix is the only digest form the corpus accepts.
const DigestPrefix = "sha256:"

// Canonical renders decoded JSON in RFC 8785 form for the value subset corpus
// documents use. Fractional and exponent numbers are rejected rather than
// formatted: no corpus document needs them, and refusing them keeps this
// encoder free of ECMAScript double-formatting rules.
func Canonical(value any) ([]byte, error) {
	buffer := &bytes.Buffer{}
	if err := writeValue(buffer, value); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Digest returns the SHA-256 digest of the canonical form of value.
func Digest(value any) (string, error) {
	canonical, err := Canonical(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return DigestPrefix + hex.EncodeToString(sum[:]), nil
}

// DecodeJSON decodes exactly one JSON document, keeping numbers exact so that
// validation and canonicalization observe the same value.
func DecodeJSON(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	value, err := decodeValue(decoder)
	if err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("unexpected content after the first JSON document")
		}
		return nil, fmt.Errorf("decode trailing content: %w", err)
	}
	return value, nil
}

func decodeValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return token, nil
	}
	switch delimiter {
	case '[':
		items := []any{}
		for decoder.More() {
			item, err := decodeValue(decoder)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		if _, err := decoder.Token(); err != nil {
			return nil, err
		}
		return items, nil
	case '{':
		members := map[string]any{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return nil, err
			}
			key, ok := keyToken.(string)
			if !ok {
				return nil, fmt.Errorf("object key has type %T", keyToken)
			}
			if _, exists := members[key]; exists {
				return nil, fmt.Errorf("duplicate object key %q", key)
			}
			value, err := decodeValue(decoder)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", key, err)
			}
			members[key] = value
		}
		if _, err := decoder.Token(); err != nil {
			return nil, err
		}
		return members, nil
	default:
		return nil, fmt.Errorf("unexpected delimiter %q", delimiter)
	}
}

func writeValue(buffer *bytes.Buffer, value any) error {
	switch typed := value.(type) {
	case nil:
		buffer.WriteString("null")
	case bool:
		buffer.WriteString(strconv.FormatBool(typed))
	case string:
		writeString(buffer, typed)
	case json.Number:
		return writeNumber(buffer, typed)
	case []any:
		return writeArray(buffer, typed)
	case map[string]any:
		return writeObject(buffer, typed)
	default:
		return fmt.Errorf("cannot canonicalize JSON value of type %T", value)
	}
	return nil
}

func writeArray(buffer *bytes.Buffer, items []any) error {
	buffer.WriteByte('[')
	for index, item := range items {
		if index > 0 {
			buffer.WriteByte(',')
		}
		if err := writeValue(buffer, item); err != nil {
			return fmt.Errorf("[%d]: %w", index, err)
		}
	}
	buffer.WriteByte(']')
	return nil
}

func writeObject(buffer *bytes.Buffer, members map[string]any) error {
	buffer.WriteByte('{')
	for index, key := range sortedKeys(members) {
		if index > 0 {
			buffer.WriteByte(',')
		}
		writeString(buffer, key)
		buffer.WriteByte(':')
		if err := writeValue(buffer, members[key]); err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
	}
	buffer.WriteByte('}')
	return nil
}

func writeNumber(buffer *bytes.Buffer, number json.Number) error {
	integer, err := strconv.ParseInt(number.String(), 10, 64)
	if err != nil {
		return fmt.Errorf("number %s is not a 64-bit integer; corpus documents may not use fractional or exponent numbers", number)
	}
	buffer.WriteString(strconv.FormatInt(integer, 10))
	return nil
}

func writeString(buffer *bytes.Buffer, value string) {
	buffer.WriteByte('"')
	for _, character := range value {
		switch character {
		case '"':
			buffer.WriteString(`\"`)
		case '\\':
			buffer.WriteString(`\\`)
		case '\b':
			buffer.WriteString(`\b`)
		case '\f':
			buffer.WriteString(`\f`)
		case '\n':
			buffer.WriteString(`\n`)
		case '\r':
			buffer.WriteString(`\r`)
		case '\t':
			buffer.WriteString(`\t`)
		default:
			if character < 0x20 {
				fmt.Fprintf(buffer, `\u%04x`, character)
				continue
			}
			buffer.WriteRune(character)
		}
	}
	buffer.WriteByte('"')
}

// sortedKeys orders member names by UTF-16 code unit, as RFC 8785 requires.
func sortedKeys(members map[string]any) []string {
	keys := make([]string, 0, len(members))
	for key := range members {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(first, second int) bool {
		return slices.Compare(utf16.Encode([]rune(keys[first])), utf16.Encode([]rune(keys[second]))) < 0
	})
	return keys
}

// IsDigest reports whether text is a well-formed corpus digest.
func IsDigest(text string) bool {
	if !strings.HasPrefix(text, DigestPrefix) {
		return false
	}
	raw, err := hex.DecodeString(strings.TrimPrefix(text, DigestPrefix))
	return err == nil && len(raw) == sha256.Size
}
