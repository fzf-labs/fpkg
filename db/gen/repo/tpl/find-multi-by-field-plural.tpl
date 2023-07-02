// FindMultiBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}查询多条数据
func (r *{{.lowerTableName}}Repo) FindMultiBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) ([]*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	dao := {{.lowerTableName}}_dao.Use(r.db).{{.upperTableName}}
	result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In({{.lowerFieldPlural}}...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}