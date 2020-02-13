package terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringTable(t *testing.T) {
	assert := assert.New(t)

	table := NewStringTable([]string{"foo", "bar"})
	table.Add("key", "value")
	s := table.String()
	assert.Contains(s, "foo")
	assert.Contains(s, "bar")
	assert.Contains(s, "key")
	assert.Contains(s, "value")
}
