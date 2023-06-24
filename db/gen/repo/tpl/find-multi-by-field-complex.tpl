func (u *{{.lowerTableName}}Repo) FindMultiBy{{.lowerFieldComplex}}(ctx context.Context, {{.lowerFields}} []{{.dataType}}) ([]*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	result, err := {{.lowerTableName}}Dao.WithContext(ctx).Where({{.lowerTableName}}Dao.Status.In(status...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}