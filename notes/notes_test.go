package notes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoteOne(t *testing.T) {
	one := NoteOne()
	assert.True(t, true, one != "")
}
