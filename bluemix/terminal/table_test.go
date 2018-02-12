package terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateStringWidth(t *testing.T) {
	for _, tc := range []struct {
		name   string
		input  string
		length int
	}{
		{name: "ascii only", input: "aB-*&%123", length: 9},
		{name: "half width unicode", input: "✣צ♥", length: 3},
		{name: "full width unicode", input: "你好", length: 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.length, calculateStringWidth(tc.input))
		})
	}
}
