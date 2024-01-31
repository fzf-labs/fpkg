package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuote(t *testing.T) {
	assert.Equal(t, Quote("abc"), "\"abc\"")
}

func TestAddSlashes(t *testing.T) {
	assert.Equal(t, AddSlashes(`{"key": 123}`), `{\"key\": 123}`)
}

func TestStripSlashes(t *testing.T) {
	assert.Equal(t, StripSlashes(`{\"key\": 123}`), `{"key": 123}`)
}

func TestUpperEnglishWord(t *testing.T) {
	assert.Equal(t, UpperEnglishWord(`wo...shi...中文`), `Wo...Shi...中文`)
}

func TestSnakeCase(t *testing.T) {
	assert.Equal(t, SnakeCase("rangePrice"), "range_price")
}

func TestSubstrReturnLeft(t *testing.T) {
	left, b := SubstrReturnLeft("AppRangePrice", "Range")
	assert.Equal(t, left, "App")
	assert.Equal(t, b, true)
}

func TestSubstrReturnRight(t *testing.T) {
	left, b := SubstrReturnRight("AppRangePrice", "Range")
	assert.Equal(t, left, "Price")
	assert.Equal(t, b, true)
}

func TestCamelCase(t *testing.T) {
	camelCase := CamelCase("App-Range-Price")
	assert.Equal(t, "appRangePrice", camelCase)
}
