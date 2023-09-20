// DeleteMultiBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}删除多条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteMultiBy{{.upperFieldPlural}}Tx(ctx context.Context,tx *{{.lowerDBName}}_dao.Query, {{.lowerFieldPlural}} []{{.dataType}}) error {
	dao := tx.{{.upperTableName}}
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
