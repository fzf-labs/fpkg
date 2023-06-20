func (u *{{.lowerTableName}}Repo) UpdateOne(ctx context.Context, data *{{.lowerTableName}}_model.{{.upperTableName}}) error {
	{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	_, err := {{.lowerTableName}}Dao.WithContext(ctx).Where({{.lowerTableName}}Dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	cache := CacheById.NewSingleKey(u.redis)
	err = cache.SingleCacheDel(ctx, data.ID)
	if err != nil {
		return err
	}
	return nil
}