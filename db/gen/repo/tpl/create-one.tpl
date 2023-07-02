func (u *{{.lowerTableName}}Repo) CreateOne(ctx context.Context, data *{{.lowerTableName}}_model.{{.upperTableName}}) error {
	dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}