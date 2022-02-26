package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilderGetDestFilename(t *testing.T) {
	b := &Builder{}

	filename := b.getDestFilename("__entity-package__/__entity-file__.go.tmpl", map[string]string{
		"entity-package": "useractivate",
		"entity-file":    "user_activate",
	})

	assert.Equal(t, "useractivate/user_activate.go", filename)
}
