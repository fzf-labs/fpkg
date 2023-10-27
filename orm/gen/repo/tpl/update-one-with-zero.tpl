// UpdateOneWithZero 更新一条数据,包含零值
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpdateOneWithZero(ctx context.Context, data *{{.dbName}}_model.{{.upperTableName}}) error {
	dao := {{.dbName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq(data.{{.upperField}})).Select(dao.ALL.WithTable("")).Updates(data)
	if err != nil {
		return err
	}
	err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, []*{{.dbName}}_model.{{.upperTableName}}{data})
    if err != nil {
    	return err
    }
	return nil
}