package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColumnNameToPropertyName(t *testing.T) {
	for _, item := range []struct {
		value    string
		expected string
	}{
		{
			value:    "id",
			expected: "ID",
		},
		{
			value:    "user_id",
			expected: "UserID",
		},
	} {
		assert.Equal(t, item.expected, ColumnNameToPropertyName(item.value))
	}
}
