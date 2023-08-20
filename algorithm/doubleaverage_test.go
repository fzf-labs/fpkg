package algorithm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoubleAverage(t *testing.T) {
	doubleAverage := DoubleAverage(100, 3)
	fmt.Println(doubleAverage)
	assert.True(t, len(doubleAverage) == 3)
}
