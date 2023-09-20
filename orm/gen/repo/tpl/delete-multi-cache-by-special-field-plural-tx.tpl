// DeleteMultiCacheBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}删除多条数据并清理缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteMultiCacheBy{{.upperFieldPlural}}Tx(ctx context.Context,tx *{{.lowerDBName}}_dao.Query, {{.lowerFieldPlural}} []{{.dataType}}) error {
	dao := tx.{{.upperTableName}}
    parameters := make([]driver.Valuer, len({{.lowerFieldPlural}}))
    for k, v := range {{.lowerFieldPlural}} {
        parameters[k] = driver.Valuer(v)
    }
	result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In(parameters...)).Find()
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.{{.upperField}}.In(parameters...)).Delete()
	if err != nil {
		return err
	}
	err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, result)
	if err != nil {
		return err
	}
	return nil
}
