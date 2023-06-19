package repo

import (
	"fmt"
	"testing"
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
}
