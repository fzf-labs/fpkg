// DeleteUniqueIndexCache 删除唯一索引存在的缓存
func (r *{{.lowerTableName}}Repo) DeleteUniqueIndexCache(ctx context.Context, data []*{{.lowerTableName}}_model.{{.upperTableName}}) error {
	var err error
    {{.singleCache}}
	for _, v := range data {
        {{.singleCacheDel}}
	}
	return nil
}