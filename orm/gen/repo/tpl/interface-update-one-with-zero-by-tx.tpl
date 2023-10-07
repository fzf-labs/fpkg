// UpdateOneWithZero 更新一条数据,包含零值(事务)
UpdateOneWithZeroByTx(ctx context.Context, tx *{{.lowerDBName}}_dao.Query, data *{{.lowerDBName}}_model.{{.upperTableName}}) error