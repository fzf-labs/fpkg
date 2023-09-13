// UpdateOne 更新一条数据(事务)
UpdateOneByTx(ctx context.Context, tx *{{.lowerDBName}}_dao.Query, data *{{.lowerDBName}}_model.{{.upperTableName}}) error