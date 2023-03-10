package algorithm

import (
	"fmt"
	"testing"
)

func TestDoubleAverage(t *testing.T) {
	doubleAverage := DoubleAverage(100, 3)
	fmt.Println(doubleAverage)
}
