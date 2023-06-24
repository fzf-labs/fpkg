func (u *{{.lowerTableName}}Repo) FindOneCacheBy{{.upperFields}}(ctx context.Context, {{.fieldsIn}}) (*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	resp := new({{.lowerTableName}}_model.{{.upperTableName}})
	cache := CacheBy{{.upperFields}}.NewSingleKey(u.redis)
	cacheValue, err := cache.SingleCache(ctx, {{.cacheKeyIn}}, func() (string, error) {
		{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
		result, err := {{.lowerTableName}}Dao.WithContext(ctx).Where({{.whereIn}}).First()
		if err != nil && err != gorm.ErrRecordNotFound {
			return "", err
		}
		marshal, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(cacheValue), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}