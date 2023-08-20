package iplocation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ipLocation_SearchLocation(t *testing.T) {
	location, err := NewIPLocation("./ip2region.xdb")
	if err != nil {
		return
	}
	searchLocation, err := location.SearchLocation("113.87.118.1")
	if err != nil {
		return
	}
	fmt.Println(searchLocation)
	assert.Equal(t, nil, err)
}

func Test_ipLocation_SearchBatchLocation(t *testing.T) {
	location, err := NewIPLocation("./ip2region.xdb")
	if err != nil {
		return
	}
	searchLocation, err := location.SearchBatchLocation([]string{"113.87.118.1", "113.87.118.2"})
	if err != nil {
		return
	}
	fmt.Println(searchLocation)
	assert.Equal(t, nil, err)
}
