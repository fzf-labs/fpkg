// UpsertOneByTx Upsert一条数据(事务)
UpsertOneByTx(ctx context.Context, tx *{{.lowerDBName}}_dao.Query, data *{{.lowerDBName}}_model.{{.upperTableName}}) error