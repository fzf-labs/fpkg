package page

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginator(t *testing.T) {
	paginator := Paginator(5, 10, 998)
	assert.Equal(t, paginator, &Page{
		Page:      5,
		PageSize:  10,
		Total:     998,
		PrevPage:  4,
		NextPage:  6,
		TotalPage: 100,
		Limit:     10,
		Offset:    40,
	})
}
