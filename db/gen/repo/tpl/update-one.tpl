// UpdateOne 更新一条数据
func (r *{{.lowerTableName}}Repo) UpdateOne(ctx context.Context, data *{{.lowerTableName}}_model.{{.upperTableName}}) error {
	dao := {{.lowerTableName}}_dao.Use(r.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	cache := CacheByID.NewSingleKey(r.redis)
	err = cache.SingleCacheDel(ctx, data.ID)
	if err != nil {
		return err
	}
	err = r.DeleteUniqueIndexCache(ctx, []*{{.lowerTableName}}_model.{{.upperTableName}}{data})
    if err != nil {
    	return err
    }
	return nil
}