// FindOneCacheBy{{.upperFields}} 根据{{.upperFields}}查询一条数据并设置缓存
func (r *{{.lowerTableName}}Repo) FindOneCacheBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) (*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	resp := new({{.lowerTableName}}_model.{{.upperTableName}})
	cache := CacheBy{{.upperFields}}.NewSingleKey(r.redis)
	cacheValue, err := cache.SingleCache(ctx, cache.BuildKey( {{.fieldsJoin}}), func() (string, error) {
		dao := {{.lowerTableName}}_dao.Use(r.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where({{.whereFields}}).First()
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