func (u *{{.lowerTableName}}Repo) FindMultiBy{{.upperFields}}(ctx context.Context, {{.fieldsIn}}) ([]*{{.lowerTableName}}_model.{{.upperTableName}}, error) {
	{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	result, err := {{.lowerTableName}}Dao.WithContext(ctx).Where({{.whereIn}}).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}
