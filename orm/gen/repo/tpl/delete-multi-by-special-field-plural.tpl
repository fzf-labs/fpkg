// DeleteMultiBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}删除多条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteMultiBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) error {
	dao := {{.dbName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
    parameters := make([]driver.Valuer, len({{.lowerFieldPlural}}))
    for k, v := range {{.lowerFieldPlural}} {
        parameters[k] = driver.Valuer(v)
    }
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In(parameters...)).Delete()
	if err != nil {
		return err
	}
	return nil
}
