// UpsertOneByTx Upsert一条数据(事务)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpsertOneByTx(ctx context.Context, tx *{{.dbName}}_dao.Query, data *{{.dbName}}_model.{{.upperTableName}}) error {
	dao := tx.{{.upperTableName}}
	err := dao.WithContext(ctx).Save(data)
	if err != nil {
		return err
	}
	return nil
}