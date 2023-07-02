func (u *{{.lowerTableName}}Repo) DeleteMultiBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) error {
	dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In({{.lowerFieldPlural}}...)).Delete()
	if err != nil {
		return err
	}
	cacheKey := CacheBy{{.lowerField}}.NewBatchKey(u.redis)
	err = cacheKey.BatchKeyCacheDel(ctx, {{.lowerFieldPlural}})
	if err != nil {
		return err
	}
	return nil
}
