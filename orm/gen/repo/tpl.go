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

//go:embed tpl/var-cache-keys.tpl
var VarCacheKeys string

//go:embed tpl/var-cache-del-key.tpl
var VarCacheDelKey string

//go:embed tpl/types.tpl
var Types string

//go:embed tpl/new.tpl
var New string

//go:embed tpl/interface-create-one.tpl
var InterfaceCreateOne string

//go:embed tpl/create-one.tpl
var CreateOne string

//go:embed tpl/interface-upsert-one.tpl
var InterfaceUpsertOne string

//go:embed tpl/upsert-one.tpl
var UpsertOne string

//go:embed tpl/interface-create-batch.tpl
var InterfaceCreateBatch string

//go:embed tpl/create-batch.tpl
var CreateBatch string

//go:embed tpl/interface-delete-multi-by-field-plural.tpl
var InterfaceDeleteMultiByFieldPlural string

//go:embed tpl/delete-multi-by-field-plural.tpl
var DeleteMultiByFieldPlural string

//go:embed tpl/interface-delete-multi-cache-by-field-plural.tpl
var InterfaceDeleteMultiCacheByFieldPlural string

//go:embed tpl/delete-multi-cache-by-field-plural.tpl
var DeleteMultiCacheByFieldPlural string

//go:embed tpl/interface-delete-one-by-field.tpl
var InterfaceDeleteOneByField string

//go:embed tpl/delete-one-by-field.tpl
var DeleteOneByField string

//go:embed tpl/interface-delete-one-cache-by-field.tpl
var InterfaceDeleteOneCacheByField string

//go:embed tpl/delete-one-cache-by-field.tpl
var DeleteOneCacheByField string

//go:embed tpl/interface-delete-one-by-fields.tpl
var InterfaceDeleteOneByFields string

//go:embed tpl/delete-one-by-fields.tpl
var DeleteOneByFields string

//go:embed tpl/interface-delete-one-cache-by-fields.tpl
var InterfaceDeleteOneCacheByFields string

//go:embed tpl/delete-one-cache-by-fields.tpl
var DeleteOneCacheByFields string

//go:embed tpl/interface-find-multi-by-field.tpl
var InterfaceFindMultiByField string

//go:embed tpl/find-multi-by-field.tpl
var FindMultiByField string

//go:embed tpl/interface-find-multi-by-fields.tpl
var InterfaceFindMultiByFields string

//go:embed tpl/find-multi-by-fields.tpl
var FindMultiByFields string

//go:embed tpl/interface-find-multi-by-field-plural.tpl
var InterfaceFindMultiByFieldPlural string

//go:embed tpl/find-multi-by-field-plural.tpl
var FindMultiByFieldPlural string

//go:embed tpl/interface-find-multi-cache-by-field-plural.tpl
var InterfaceFindMultiCacheByFieldPlural string

//go:embed tpl/find-multi-cache-by-field-plural.tpl
var FindMultiCacheByFieldPlural string

//go:embed tpl/interface-find-one-by-field.tpl
var InterfaceFindOneByField string

//go:embed tpl/find-one-by-field.tpl
var FindOneByField string

//go:embed tpl/interface-find-one-cache-by-field.tpl
var InterfaceFindOneCacheByField string

//go:embed tpl/find-one-cache-by-field.tpl
var FindOneCacheByField string

//go:embed tpl/interface-find-one-by-fields.tpl
var InterfaceFindOneByFields string

//go:embed tpl/find-one-by-fields.tpl
var FindOneByFields string

//go:embed tpl/interface-find-one-cache-by-fields.tpl
var InterfaceFindOneCacheByFields string

//go:embed tpl/find-one-cache-by-fields.tpl
var FindOneCacheByFields string

//go:embed tpl/interface-find-multi-by-paginator.tpl
var InterfaceFindMultiByPaginator string

//go:embed tpl/find-multi-by-paginator.tpl
var FindMultiByPaginator string

//go:embed tpl/interface-update-one.tpl
var InterfaceUpdateOne string

//go:embed tpl/update-one.tpl
var UpdateOne string

//go:embed tpl/interface-delete-unique-index-cache.tpl
var InterfaceDeleteUniqueIndexCache string

//go:embed tpl/delete-unique-index-cache.tpl
var DeleteUniqueIndexCache string

//go:embed tpl/interface-create-one-by-tx.tpl
var InterfaceCreateOneByTx string

//go:embed tpl/create-one-by-tx.tpl
var CreateOneByTx string

//go:embed tpl/interface-upsert-one-by-tx.tpl
var InterfaceUpsertOneByTx string

//go:embed tpl/upsert-one-by-tx.tpl
var UpsertOneByTx string

//go:embed tpl/interface-update-one-by-tx.tpl
var InterfaceUpdateOneByTx string

//go:embed tpl/update-one-by-tx.tpl
var UpdateOneByTx string

//go:embed tpl/interface-delete-multi-by-field-plural-tx.tpl
var InterfaceDeleteMultiByFieldPluralTx string

//go:embed tpl/delete-multi-by-field-plural-tx.tpl
var DeleteMultiByFieldPluralTx string

//go:embed tpl/interface-delete-multi-cache-by-field-plural-tx.tpl
var InterfaceDeleteMultiCacheByFieldPluralTx string

//go:embed tpl/delete-multi-cache-by-field-plural-tx.tpl
var DeleteMultiCacheByFieldPluralTx string

//go:embed tpl/interface-delete-one-by-field-tx.tpl
var InterfaceDeleteOneByFieldTx string

//go:embed tpl/delete-one-by-field-tx.tpl
var DeleteOneByFieldTx string

//go:embed tpl/interface-delete-one-cache-by-field-tx.tpl
var InterfaceDeleteOneCacheByFieldTx string

//go:embed tpl/delete-one-cache-by-field-tx.tpl
var DeleteOneCacheByFieldTx string

//go:embed tpl/interface-delete-one-by-fields-tx.tpl
var InterfaceDeleteOneByFieldsTx string

//go:embed tpl/delete-one-by-fields-tx.tpl
var DeleteOneByFieldsTx string

//go:embed tpl/interface-delete-one-cache-by-fields-tx.tpl
var InterfaceDeleteOneCacheByFieldsTx string

//go:embed tpl/delete-one-cache-by-fields-tx.tpl
var DeleteOneCacheByFieldsTx string
