// CreateOneByTx 创建一条数据(事务)
CreateOneByTx(ctx context.Context, tx *{{.dbName}}_dao.Query, data *{{.dbName}}_model.{{.upperTableName}}) error