// FindOneCacheBy{{.upperFields}} 根据{{.upperFields}}查询一条数据并设置缓存
func (r *{{.upperTableName}}Repo) FindOneCacheBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) (*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
	resp := new({{.lowerDBName}}_model.{{.upperTableName}})
	key := r.cache.Key(ctx, cache{{.upperTableName}}By{{.upperFields}}Prefix, {{.lowerFieldsJoin}})
	keys := []string{key}
	cacheValue, err := r.cache.Take(ctx, keys, func() (map[string]string, error) {
		dao := {{.lowerDBName}}_dao.Use(r.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where({{.whereFields}}).First()
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


