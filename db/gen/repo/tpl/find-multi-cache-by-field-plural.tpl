// FindMultiCacheBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}查询多条数据并设置缓存
func (r *{{.upperTableName}}Repo) FindMultiCacheBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) ([]*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
	resp := make([]*{{.lowerDBName}}_model.{{.upperTableName}}, 0)
	keys := make([]string, 0)
	for _, v := range {{.lowerFieldPlural}} {
		keys = append(keys, r.cache.Key(ctx, cache{{.upperTableName}}By{{.upperField}}Prefix, v))
	}
	cacheValue, err := r.cache.Take(ctx, keys, func() (map[string]string, error) {
		dao := {{.lowerDBName}}_dao.Use(r.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In({{.lowerFieldPlural}}...)).Find()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		value := make(map[string]string)
		for _, v := range keys {
			value[v] = ""
		}
		for _, v := range result {
			marshal, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			value[r.cache.Key(ctx, cache{{.upperTableName}}By{{.upperField}}Prefix, v.{{.upperField}})] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new({{.lowerDBName}}_model.{{.upperTableName}})
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}