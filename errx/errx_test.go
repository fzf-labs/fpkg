package errx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := callersStack(2, 1)
	str := s.String()
	assert.True(t, len(str) > 0)
}
