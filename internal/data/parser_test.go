package data

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidRecord(t *testing.T) {
	tests := []struct {
		record   []string
		expected bool
	}{
		{[]string{"1", `{"type":"Feature"}`, "2025-02-09T15:04:05Z"}, true},
		{[]string{"1", "", "2025-02-09T15:04:05Z"}, false},
		{[]string{"1", `{"type":"Feature"}`, "invalid_timestamp"}, false},
		{[]string{`{"type":"Feature"}`, "2025-02-09T15:04:05Z"}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Testing record %v", tt.record), func(t *testing.T) {
			result := isValidRecord(tt.record)
			assert.Equal(t, tt.expected, result)
		})
	}
}
