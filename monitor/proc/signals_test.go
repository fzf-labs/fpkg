package proc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDone(t *testing.T) {
	select {
	case <-Done():
		assert.Fail(t, "should run")
	default:
	}
	assert.NotNil(t, Done())
}

func TestRun(t *testing.T) {
	// kill -usr1 xxx
	// kill -usr2 xxx
	Monitor()
	assert.Equal(t, nil, nil)
}
