// DeleteUniqueIndexCache 删除唯一索引存在的缓存
func (r *{{.upperTableName}}Repo) DeleteUniqueIndexCache(ctx context.Context, data []*{{.lowerDBName}}_model.{{.upperTableName}}) error {
	keys := make([]string, 0)
	for _, v := range data {
	  {{.varCacheDelKeys}}
	}
	err := r.cache.DelBatch(ctx, keys)
	if err != nil {
		return err
	}
	return nil
}