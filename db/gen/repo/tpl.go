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

//go:embed tpl/var-single-cache.tpl
var VarSingleCache string

//go:embed tpl/var-single-cache-del.tpl
var VarSingleCacheDel string

//go:embed tpl/types.tpl
var Types string

//go:embed tpl/new.tpl
var New string

//go:embed tpl/interface-create-one.tpl
var InterfaceCreateOne string

//go:embed tpl/create-one.tpl
var CreateOne string

//go:embed tpl/interface-delete-multi-cache-by-field-plural.tpl
var InterfaceDeleteMultiCacheByFieldPlural string

//go:embed tpl/delete-multi-cache-by-field-plural.tpl
var DeleteMultiCacheByFieldPlural string

//go:embed tpl/interface-delete-one-cache-by-field.tpl
var InterfaceDeleteOneCacheByField string

//go:embed tpl/delete-one-cache-by-field.tpl
var DeleteOneCacheByField string

//go:embed tpl/interface-delete-one-cache-by-fields.tpl
var InterfaceDeleteOneCacheByFields string

//go:embed tpl/delete-one-cache-by-fields.tpl
var DeleteOneCacheByFields string

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

//go:embed tpl/interface-delete-unique-index-cache.tpl
var InterfaceDeleteUniqueIndexCache string

//go:embed tpl/delete-unique-index-cache.tpl
var DeleteUniqueIndexCache string
