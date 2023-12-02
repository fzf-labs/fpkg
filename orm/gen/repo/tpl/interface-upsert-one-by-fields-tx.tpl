// UpsertOneByFieldsTx Upsert一条数据，根据fields字段(事务)
UpsertOneByFieldsTx(ctx context.Context,tx *{{.dbName}}_dao.Query, data *{{.dbName}}_model.{{.upperTableName}},fields []string) error