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

//go:embed tpl/new.tpl
var New string

//go:embed tpl/interface-create-one.tpl
var InterfaceCreateOne string

//go:embed tpl/create-one.tpl
var CreateOne string

//go:embed tpl/interface-delete-multi-by-field-plural.tpl
var InterfaceDeleteMultiByFieldPlural string

//go:embed tpl/delete-multi-by-field-plural.tpl
var DeleteMultiByFieldPlural string

//go:embed tpl/interface-delete-one-by-field.tpl
var InterfaceDeleteOneByField string

//go:embed tpl/delete-one-by-field.tpl
var DeleteOneByField string

//go:embed tpl/interface-find-multi-by-field-plural.tpl
var InterfaceFindMultiByFieldPlural string

//go:embed tpl/find-multi-by-field-plural.tpl
var FindMultiByFieldPlural string

//go:embed tpl/interface-find-multi-by-fields.tpl
var InterfaceFindMultiByFields string

//go:embed tpl/find-multi-by-fields.tpl
var FindMultiByFields string

//go:embed tpl/interface-find-multi-cache-by-field-plural.tpl
var InterfaceFindMultiCacheByFieldPlural string

//go:embed tpl/find-multi-cache-by-field-plural.tpl
var FindMultiCacheByFieldPlural string

//go:embed tpl/interface-find-one-cache-by-field.tpl
var InterfaceFindOneCacheByField string

//go:embed tpl/find-one-cache-by-field.tpl
var FindOneCacheByField string

//go:embed tpl/interface-find-one-cache-by-fields.tpl
var InterfaceFindOneCacheByFields string

//go:embed tpl/find-one-cache-by-fields.tpl
var FindOneCacheByFields string

//go:embed tpl/interface-update-one.tpl
var InterfaceUpdateOne string

//go:embed tpl/update-one.tpl
var UpdateOne string

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
	lowerField       string //字段小写
	upperField       string //大写字段
	lowerFieldPlural string //字段小写复数
	upperFieldPlural string //字段大写复数
	dataType         string //字段类型
}
