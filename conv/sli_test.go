package conv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliToMap(t *testing.T) {
	sli1 := []string{"1", "2", "3"}

	sliToMap := SliToMap(sli1, func(s string) (string, string) {
		return s, s
	})
	fmt.Println(sliToMap)
	assert.True(t, len(sliToMap) > 0)
}
