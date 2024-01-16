// FindOneCacheBy{{.upperField}} 根据{{.lowerField}}查询一条数据并设置缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindOneCacheBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.dbName}}_model.{{.upperTableName}}, error) {
	resp := new({{.dbName}}_model.{{.upperTableName}})
	cacheKey := {{.firstTableChar}}.cache.Key( cache{{.upperTableName}}By{{.upperField}}Prefix, {{.lowerField}})
	cacheValue, err := {{.firstTableChar}}.cache.Fetch(ctx, cacheKey, func() (string, error) {
		dao := {{.dbName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where({{.whereField}}).First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
	if cacheValue != "" {
		err = json.Unmarshal([]byte(cacheValue), resp)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}