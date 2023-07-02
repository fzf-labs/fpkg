// DeleteMultiCacheBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}删除多条数据并清理缓存
func (r *{{.lowerTableName}}Repo) DeleteMultiCacheBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) error {
	dao := {{.lowerTableName}}_dao.Use(r.db).{{.upperTableName}}
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
	err = r.DeleteUniqueIndexCache(ctx, list)
	if err != nil {
		return err
	}
	return nil
}
