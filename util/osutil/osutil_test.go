package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsWindows(t *testing.T) {
	assert.True(t, false, IsWindows())
}
