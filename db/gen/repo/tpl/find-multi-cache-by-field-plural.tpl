// FindMultiCacheBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}查询多条数据并设置缓存
func (r *{{.upperTableName}}Repo) FindMultiCacheBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) ([]*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
	resp := make([]*{{.lowerDBName}}_model.{{.upperTableName}}, 0)
	keys := make([]string, 0)
	keyToParam := make(map[string]{{.dataType}})
	for _, v := range {{.lowerFieldPlural}} {
	    key := r.cache.Key(ctx, cache{{.upperTableName}}By{{.upperField}}Prefix, v)
		keys = append(keys,key)
		keyToParam[key] = v
	}
	cacheValue, err := r.cache.FetchBatch(ctx, keys, func(miss []string) (map[string]string, error) {
        params := make([]{{.dataType}},0)
        for _, v := range miss {
            params = append(params, keyToParam[v])
        }
		dao := {{.lowerDBName}}_dao.Use(r.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In(params...)).Find()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		value := make(map[string]string)
		for _, v := range miss {
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