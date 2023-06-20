func (u *{{.lowerTableName}}Repo) FindMultiBy{{.upperFields}}(ctx context.Context, {{.lowerFields}} []{{.dataType}}) ([]*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	resp := make([]*{{.lowerTableName}}_model.{{.upperTableName}}, 0)
	cacheKey := CacheBy{{.lowerField}}.NewBatchKey(u.redis)
	cacheValue, err := cacheKey.BatchKeyCache(ctx, {{.lowerFields}}, func() (map[string]string, error) {
		{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
		result, err := {{.lowerTableName}}Dao.WithContext(ctx).Where({{.lowerTableName}}Dao.{{.upperField}}.In({{.lowerFields}}...)).Find()
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