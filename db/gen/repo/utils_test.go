package repo

import (
	"fmt"
	"testing"

	"github.com/jinzhu/inflection"
)

func TestUpperName(t *testing.T) {
	fmt.Println(UpperName("id"))
	fmt.Println(inflection.Plural(UpperName("id")))
	fmt.Println(inflection.Plural(LowerName("id")))
}
