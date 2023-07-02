func (u *{{.lowerTableName}}Repo) FindMultiBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) ([]*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	result, err := dao.WithContext(ctx).Where({{.whereFields}}).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}
