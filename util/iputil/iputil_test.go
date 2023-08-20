package iputil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicIP(t *testing.T) {
	ip := GetPublicIP()
	assert.True(t, ip != "")
}

func TestGetPublicIPByHTTP(t *testing.T) {
	ip, err := GetPublicIPByHTTP()
	if err != nil {
		return
	}
	assert.True(t, ip != "")
	assert.Equal(t, nil, err)
}
