// FindOneBy{{.upperField}} 根据{{.lowerField}}查询一条数据并设置缓存
func (r *{{.lowerTableName}}Repo) FindOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	resp := new({{.lowerTableName}}_model.{{.upperTableName}})
	cache := CacheBy{{.upperField}}.NewSingleKey(r.redis)
	cacheValue, err := cache.SingleCache(ctx, {{.lowerField}} , func() (string, error) {
		dao := {{.lowerTableName}}_dao.Use(r.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).First()
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