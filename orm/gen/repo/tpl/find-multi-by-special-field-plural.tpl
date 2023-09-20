// FindMultiBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}查询多条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) ([]*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
    parameters := make([]driver.Valuer, len({{.lowerFieldPlural}}))
    for k, v := range {{.lowerFieldPlural}} {
        parameters[k] = driver.Valuer(v)
    }
	result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In(parameters...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}