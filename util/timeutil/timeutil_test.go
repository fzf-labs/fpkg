package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func TestNowUnix(t *testing.T) {
	fmt.Println(time.Now().String())
	fmt.Println(time.Now().GoString())
}
