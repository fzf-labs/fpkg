package proc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfile(t *testing.T) {
	profiler := StartProfile()
	// start again should not work
	assert.NotNil(t, StartProfile())
	profiler.Stop()
	// stop twice
	profiler.Stop()
	assert.True(t, strings.Contains("", ".pprof"))
}
