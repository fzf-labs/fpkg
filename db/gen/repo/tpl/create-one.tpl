func (u *{{.lowerTableName}}Repo) CreateOne(ctx context.Context, data *{{.lowerTableName}}_model.User) error {
	{{.lowerTableName}}Dao := {{.lowerTableName}}_dao.Use(u.db).{{.upperTableName}}
	err := {{.lowerTableName}}Dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}