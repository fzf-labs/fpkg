// FindOneCacheBy{{.upperField}} 根据{{.lowerField}}查询一条数据并设置缓存
func (r *{{.upperTableName}}Repo) FindOneCacheBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
	resp := new({{.lowerDBName}}_model.{{.upperTableName}})
	key := r.cache.Key(ctx, cache{{.upperTableName}}By{{.upperField}}Prefix, {{.lowerField}})
	keys := []string{key}
	cacheValue, err := r.cache.Take(ctx, keys, func() (map[string]string, error) {
		dao := {{.lowerDBName}}_dao.Use(r.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		value := make(map[string]string)
        marshal, err := json.Marshal(result)
        if err != nil {
            return nil, err
        }
        value[key] = string(marshal)
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(cacheValue[key]), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}