package urlutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLEncode(t *testing.T) {
	assert.Equal(t, URLEncode("www.baidu.com/?query=golang"), "www.baidu.com/?query%3Dgolang")
}

func TestURLDecode(t *testing.T) {
	assert.Equal(t, URLDecode("www.baidu.com/?query%3Dgolang"), "www.baidu.com/?query=golang")
}
