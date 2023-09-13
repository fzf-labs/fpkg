// FindOneCacheBy{{.upperField}} 根据{{.lowerField}}查询一条数据并设置缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindOneCacheBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
	resp := new({{.lowerDBName}}_model.{{.upperTableName}})
	key := {{.firstTableChar}}.cache.Key( cache{{.upperTableName}}By{{.upperField}}Prefix, {{.lowerField}})
	cacheValue, err := {{.firstTableChar}}.cache.Fetch(ctx, key, func() (string, error) {
		dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
		result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).First()
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