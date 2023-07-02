// FindMultiBy{{.upperFields}} 根据{{.upperFields}}查询多条数据
func (r *{{.lowerTableName}}Repo) FindMultiBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) ([]*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	dao := {{.lowerTableName}}_dao.Use(r.db).{{.upperTableName}}
	result, err := dao.WithContext(ctx).Where({{.whereFields}}).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}
