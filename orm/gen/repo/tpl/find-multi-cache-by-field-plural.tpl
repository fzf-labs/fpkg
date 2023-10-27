// FindMultiCacheBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}查询多条数据并设置缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiCacheBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) ([]*{{.dbName}}_model.{{.upperTableName}}, error) {
	resp := make([]*{{.dbName}}_model.{{.upperTableName}}, 0)
	cacheKeys := make([]string, 0)
	keyToParam := make(map[string]{{.dataType}})
	for _, v := range {{.lowerFieldPlural}} {
	    cacheKey := {{.firstTableChar}}.cache.Key( cache{{.upperTableName}}By{{.upperField}}Prefix, v)
		cacheKeys = append(cacheKeys,cacheKey)
		keyToParam[cacheKey] = v
	}
	cacheValue, err := {{.firstTableChar}}.cache.FetchBatch(ctx, cacheKeys, func(miss []string) (map[string]string, error) {
        parameters := make([]{{.dataType}},0)
        for _, v := range miss {
            parameters = append(parameters, keyToParam[v])
        }
		dao := {{.dbName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In(parameters...)).Find()
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
			value[{{.firstTableChar}}.cache.Key( cache{{.upperTableName}}By{{.upperField}}Prefix, v.{{.upperField}})] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new({{.dbName}}_model.{{.upperTableName}})
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}