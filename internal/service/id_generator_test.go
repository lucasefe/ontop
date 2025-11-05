package service

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateID(t *testing.T) {
	// Test basic generation
	id := GenerateID()

	// Verify length (26 characters)
	if len(id) != 26 {
		t.Errorf("Expected ID length 26, got %d: %s", len(id), id)
	}

	// Verify it only contains Crockford Base32 alphabet
	alphabet := "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	for _, char := range id {
		if !strings.ContainsRune(alphabet, char) {
			t.Errorf("ID contains invalid character: %c in %s", char, id)
		}
	}
}

func TestGenerateID_Uniqueness(t *testing.T) {
	// Generate multiple IDs and verify they're unique
	ids := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		id := GenerateID()
		if ids[id] {
			t.Errorf("Duplicate ID generated: %s", id)
		}
		ids[id] = true
	}

	if len(ids) != count {
		t.Errorf("Expected %d unique IDs, got %d", count, len(ids))
	}
}

func TestGenerateID_Sortable(t *testing.T) {
	// Generate IDs with small delays and verify they're lexicographically sorted
	ids := make([]string, 10)

	for i := 0; i < 10; i++ {
		ids[i] = GenerateID()
		time.Sleep(2 * time.Millisecond) // Small delay to ensure different timestamps
	}

	// Verify each ID is greater than or equal to the previous (lexicographic order)
	for i := 1; i < len(ids); i++ {
		if ids[i] < ids[i-1] {
			t.Errorf("IDs not in lexicographic order: %s should be >= %s", ids[i], ids[i-1])
		}
	}
}

func TestGenerateID_Format(t *testing.T) {
	id := GenerateID()

	// Test that timestamp part is encoded correctly
	// First character should represent a reasonable timestamp
	// (we're in year 2025, so first char should be '0' or '1')
	firstChar := string(id[0])
	if firstChar != "0" && firstChar != "1" {
		t.Logf("Warning: First character is %s, expected 0 or 1 for current timestamp", firstChar)
	}

	t.Logf("Generated ULID: %s", id)
}
