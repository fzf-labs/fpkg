// UpsertOneByTx Upsert一条数据(事务)
UpsertOneByTx(ctx context.Context, tx *{{.dbName}}_dao.Query, data *{{.dbName}}_model.{{.upperTableName}}) error