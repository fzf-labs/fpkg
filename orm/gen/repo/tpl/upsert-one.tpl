// UpsertOne Upsert一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpsertOne(ctx context.Context, data *{{.lowerDBName}}_model.{{.upperTableName}}) error {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	err := dao.WithContext(ctx).Save(data)
	if err != nil {
		return err
	}
	return nil
}