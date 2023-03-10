package page

import (
	"fmt"
	"testing"
)

func TestPaginator(t *testing.T) {
	paginator := Paginator(5, 10, 998)
	fmt.Println(paginator)
}
