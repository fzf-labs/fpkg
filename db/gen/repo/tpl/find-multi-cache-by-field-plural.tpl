// FindMultiBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}查询多条数据并设置缓存
func (r *{{.lowerTableName}}Repo) FindMultiBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) ([]*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	resp := make([]*{{.lowerTableName}}_model.{{.upperTableName}}, 0)
	cacheKey := CacheBy{{.upperField}}.NewBatchKey(r.redis)
	cacheValue, err := cacheKey.BatchKeyCache(ctx, {{.lowerFieldPlural}}, func() (map[string]string, error) {
		dao := {{.lowerTableName}}_dao.Use(r.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In({{.lowerFieldPlural}}...)).Find()
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		value := make(map[string]string)
		for _, v := range result {
			marshal, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			value[v.{{.upperField}}] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new({{.lowerTableName}}_model.{{.upperTableName}})
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}