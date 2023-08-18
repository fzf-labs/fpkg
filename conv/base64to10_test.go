package conv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase10To64(t *testing.T) {
	to64, err := Base10To64(29)
	fmt.Println(to64)
	fmt.Println(err)
	assert.Equal(t, "t", to64)
	assert.Equal(t, nil, err)
}

func TestBase64To10(t *testing.T) {
	to64, err := Base64To10("0t")
	assert.Equal(t, 29, int(to64))
	assert.Equal(t, nil, err)
}
