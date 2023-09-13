// UpdateOneByTx 更新一条数据(事务)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpdateOneByTx(ctx context.Context, tx *{{.lowerDBName}}_dao.Query, data *{{.lowerDBName}}_model.{{.upperTableName}}) error {
	dao := tx.{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
    err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, []*{{.lowerDBName}}_model.{{.upperTableName}}{data})
    if err != nil {
        return err
    }
	return err
}