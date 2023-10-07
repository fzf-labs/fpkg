// UpdateOneWithZero 更新一条数据,包含零值
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpdateOneWithZero(ctx context.Context, data *{{.lowerDBName}}_model.{{.upperTableName}}) error {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq(data.{{.upperField}})).Select(dao.ALL).Updates(data)
	if err != nil {
		return err
	}
	err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, []*{{.lowerDBName}}_model.{{.upperTableName}}{data})
    if err != nil {
    	return err
    }
	return nil
}