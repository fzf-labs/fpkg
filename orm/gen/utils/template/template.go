package template

import (
	"bytes"
	"go/format"
	"text/template"

	"github.com/pkg/errors"
)

// DefaultTemplate is a tool to provides the text/template operations
type DefaultTemplate struct {
	name  string
	text  string
	goFmt bool
}

// NewTemplate returns an instance of defaultTemplate
func NewTemplate(name string) *DefaultTemplate {
	return &DefaultTemplate{
		name: name,
	}
}

// Parse accepts a source template and returns defaultTemplate
func (t *DefaultTemplate) Parse(text string) *DefaultTemplate {
	t.text = text
	return t
}

// GoFmt sets the value to goFmt and marks the generated codes will be formatted or not
func (t *DefaultTemplate) GoFmt(fmt bool) *DefaultTemplate {
	t.goFmt = fmt
	return t
}

// Execute returns the codes after the template executed
func (t *DefaultTemplate) Execute(data any) (*bytes.Buffer, error) {
	tem, err := template.New(t.name).Parse(t.text)
	if err != nil {
		return nil, errors.Wrapf(err, "template parse error:%s", t.text)
	}

	buf := new(bytes.Buffer)
	if err = tem.Execute(buf, data); err != nil {
		return nil, errors.Wrapf(err, "template execute error:%s", t.text)
	}

	if !t.goFmt {
		return buf, nil
	}

	formatOutput, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "go format error:%s", buf.String())
	}

	buf.Reset()
	buf.Write(formatOutput)
	return buf, nil
}
