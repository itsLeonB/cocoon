package util_test

import (
	"testing"

	"github.com/itsLeonB/cocoon/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestGetNameFromEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected string
	}{
		{"john.doe@example.com", "John"},
		{"jane_doe123@example.com", "Jane"},
		{"1234@example.com", ""},
		{"no-symbol.com", ""},
		{"@example.com", ""},
		{"", ""},
		{"alice@example.com", "Alice"},
		{"Bob_TheBuilder99@example.com", "Bob"},
		{"_hidden.user@example.com", "Hidden"},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := util.GetNameFromEmail(tt.email)
			assert.Equal(t, tt.expected, result, "Expected result for email %q", tt.email)
		})
	}
}
