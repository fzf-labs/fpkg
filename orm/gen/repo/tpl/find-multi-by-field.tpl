// FindMultiBy{{.upperField}} 根据{{.lowerField}}查询多条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) ([]*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	result, err := dao.WithContext(ctx).Where({{.whereField}}).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}