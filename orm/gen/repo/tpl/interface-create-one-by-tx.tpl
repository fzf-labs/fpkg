// CreateOneByTx 创建一条数据(事务)
CreateOneByTx(ctx context.Context, tx *{{.lowerDBName}}_dao.Query, data *{{.lowerDBName}}_model.{{.upperTableName}}) error