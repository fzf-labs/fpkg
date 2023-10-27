// UpdateOneByTx 更新一条数据(事务)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpdateOneByTx(ctx context.Context, tx *{{.dbName}}_dao.Query, data *{{.dbName}}_model.{{.upperTableName}}) error {
	dao := tx.{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq(data.{{.upperField}})).Updates(data)
	if err != nil {
		return err
	}
    err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, []*{{.dbName}}_model.{{.upperTableName}}{data})
    if err != nil {
        return err
    }
	return err
}