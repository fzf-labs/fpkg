// DeleteMultiCacheBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}删除多条数据并清理缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteMultiCacheBy{{.upperFieldPlural}}Tx(ctx context.Context,tx *{{.lowerDBName}}_dao.Query, {{.lowerFieldPlural}} []{{.dataType}}) error {
	dao := tx.{{.upperTableName}}
	list, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In({{.lowerFieldPlural}}...)).Find()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.{{.upperField}}.In({{.lowerFieldPlural}}...)).Delete()
	if err != nil {
		return err
	}
	err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, list)
	if err != nil {
		return err
	}
	return nil
}
