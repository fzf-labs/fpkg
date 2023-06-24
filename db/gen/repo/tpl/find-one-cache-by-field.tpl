func (u *{{.lowerTableName}}Repo) FindOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	resp := new({{.lowerTableName}}_model.{{.upperTableName}})
	cache := CacheBy{{.upperField}}.NewSingleKey(u.redis)
	cacheValue, err := cache.SingleCache(ctx, {{.lowerField}} , func() (string, error) {
		{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
		result, err := {{.lowerTableName}}Dao.WithContext(ctx).Where({{.lowerTableName}}Dao.{{.upperField}}.Eq({{.lowerField}})).First()
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