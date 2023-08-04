// FindMultiBy{{.upperField}} 根据{{.lowerField}}查询多条数据
func (r *{{.upperTableName}}Repo) FindMultiBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) ([]*{{.lowerDbName}}_model.{{.upperTableName}}, error) {
	dao := {{.lowerDbName}}_dao.Use(r.db).{{.upperTableName}}
	result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}