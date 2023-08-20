package errx

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunc_FileLine(t *testing.T) {
	// str := "abc"
	caller, file, line, ok := runtime.Caller(0)
	fmt.Println(caller, file, line, ok)
	assert.Equal(t, true, ok)
}
