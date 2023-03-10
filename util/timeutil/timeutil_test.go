package timeutil

import (
	"fmt"
	"testing"
)

func TestNowUnix(t *testing.T) {
	fmt.Println(NowUnix())
}
