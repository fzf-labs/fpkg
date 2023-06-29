package repo

import (
	_ "embed"
)

//go:embed tpl/pkg.tpl
var Pkg string

//go:embed tpl/import.tpl
var Import string

//go:embed tpl/var.tpl
var Var string

//go:embed tpl/var-cache.tpl
var VarCache string

//go:embed tpl/types.tpl
var Types string

//go:embed tpl/interface-create-one.tpl
var InterfaceCreateOne string

//go:embed tpl/interface-delete-multi-by-field-complex.tpl
var InterfaceDeleteMultiByFieldComplex string

//go:embed tpl/interface-delete-one-by-field.tpl.tpl
var InterfaceDeleteOneByField string

//go:embed tpl/interface-find-multi-by-field-complex.tpl
var InterfaceFindMultiByFieldComplex string

//go:embed tpl/interface-find-multi-by-fields.tpl
var InterfaceFindMultiByFields string

//go:embed tpl/interface-find-multi-cache-by-field-complex.tpl
var InterfaceFindMultiCacheByFieldComplex string

//go:embed tpl/interface-find-one-cache-by-field.tpl
var InterfaceFindOneCacheByField string

//go:embed tpl/interface-find-one-cache-by-fields.tpl
var InterfaceFindOneCacheByFields string

//go:embed tpl/interface-update-one.tpl
var InterfaceUpdateOne string

// TableRepo 表结构
type TableRepo struct {
	dbName         string
	lowerDbName    string
	upperDbName    string
	tableName      string
	lowerTableName string
	upperTableName string
	cacheKeys      []string
	methods        []string
}

type FuncRepo struct {
	lowerField string
	upperField string
	in         string
}

type FieldRepo struct {
	lowerField        string //字段小写
	upperField        string //大写字段
	lowerFieldComplex string //字段小写复数
	upperFieldComplex string //字段大写复数
	dataType          string //字段类型
}
