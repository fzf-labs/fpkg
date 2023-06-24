func (u *{{.lowerTableName}}Repo) DeleteOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) error {
	{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	_, err := {{.lowerTableName}}Dao.WithContext(ctx).Where({{.lowerTableName}}Dao.{{.upperField}}.Eq({{.lowerField}})).Delete()
	if err != nil {
		return err
	}
	cache := CacheBy{{.upperField}}.NewSingleKey(u.redis)
	err = cache.SingleCacheDel(ctx, {{.lowerField}})
	if err != nil {
		return err
	}
	return nil
}