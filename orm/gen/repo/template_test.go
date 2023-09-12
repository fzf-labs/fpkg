package repo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplate(t *testing.T) {
	execute, err := NewTemplate("insert").Parse(Pkg).Execute(map[string]any{
		"withCache":             true,
		"upperStartCamelObject": "",
		"lowerStartCamelObject": "",
		"expression":            "",
		"expressionValues":      "",
		"keys":                  "",
		"keyValues":             "",
		"data":                  "",
	})
	if err != nil {
		return
	}
	fmt.Println(execute)
	assert.Equal(t, nil, err)
}