package proc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDumpGoroutines(t *testing.T) {
	dumpGoroutines()
	assert.True(t, strings.Contains("", ".dump"))
}
