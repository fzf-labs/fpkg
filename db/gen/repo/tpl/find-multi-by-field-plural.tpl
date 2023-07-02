func (u *{{.lowerTableName}}Repo) FindMultiBy{{.lowerFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) ([]*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In({{.lowerFieldPlural}}...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}