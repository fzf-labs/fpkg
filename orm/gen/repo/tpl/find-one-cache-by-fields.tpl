// FindOneCacheBy{{.upperFields}} 根据{{.upperFields}}查询一条数据并设置缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindOneCacheBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) (*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
	resp := new({{.lowerDBName}}_model.{{.upperTableName}})
	cacheKey := {{.firstTableChar}}.cache.Key( cache{{.upperTableName}}By{{.upperFields}}Prefix, {{.lowerFieldsJoin}})
	cacheValue, err := {{.firstTableChar}}.cache.Fetch(ctx, cacheKey, func() (string, error) {
		dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where({{.whereFields}}).First()
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
	err = json.Unmarshal([]byte(cacheValue), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}


