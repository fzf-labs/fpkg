package proto

import (
	_ "embed"
)

//go:embed tpl/syntax.tpl
var Syntax string

//go:embed tpl/package.tpl
var Package string

//go:embed tpl/import.tpl
var Import string

//go:embed tpl/option.tpl
var Option string

//go:embed tpl/service.tpl
var Service string

//go:embed tpl/message.tpl
var Message string
