func (u *{{.lowerTableName}}Repo) DeleteMultiBy{{.upperFields}}(ctx context.Context, {{.lowerFields}} []{{.dataType}}) error {
	{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	_, err := {{.lowerTableName}}Dao.WithContext(ctx).Where({{.lowerTableName}}Dao.{{.upperField}}.In({{.lowerFields}}...)).Delete()
	if err != nil {
		return err
	}
	cacheKey := CacheBy{{.lowerField}}.NewBatchKey(u.redis)
	err = cacheKey.BatchKeyCacheDel(ctx, {{.lowerFields}})
	if err != nil {
		return err
	}
	return nil
}
