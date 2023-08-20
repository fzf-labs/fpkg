package goscraper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrape(t *testing.T) {
	scrape, err := Scrape("https://www.baidu.com", 1)
	if err != nil {
		return
	}
	fmt.Println(scrape)
	assert.Equal(t, nil, err)
}
