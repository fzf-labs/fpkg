// UpdateOne 更新一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpdateOne(ctx context.Context, data *{{.lowerDBName}}_model.{{.upperTableName}}) error {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, []*{{.lowerDBName}}_model.{{.upperTableName}}{data})
    if err != nil {
    	return err
    }
	return nil
}