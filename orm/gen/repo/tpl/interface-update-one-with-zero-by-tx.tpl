// UpdateOneWithZero 更新一条数据,包含零值(事务)
UpdateOneWithZeroByTx(ctx context.Context, tx *{{.dbName}}_dao.Query, data *{{.dbName}}_model.{{.upperTableName}}) error