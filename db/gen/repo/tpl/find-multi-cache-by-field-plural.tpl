// FindMultiCacheBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}查询多条数据并设置缓存
func (r *{{.upperTableName}}Repo) FindMultiCacheBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) ([]*{{.lowerDbName}}_model.{{.upperTableName}}, error) {
	resp := make([]*{{.lowerDbName}}_model.{{.upperTableName}}, 0)
	cacheKey := Cache{{.upperTableName}}By{{.upperField}}.NewBatchKey(r.redis)
	batchKeys := make([]string,0)
	for _, v := range {{.lowerFieldPlural}} {
		batchKeys = append(batchKeys,conv.String(v))
	}
	cacheValue, err := cacheKey.BatchKeyCache(ctx, batchKeys, func() (map[string]string, error) {
		dao := {{.lowerDbName}}_dao.Use(r.db).{{.upperTableName}}
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
			value[conv.String(v.{{.upperField}})] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new({{.lowerDbName}}_model.{{.upperTableName}})
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}