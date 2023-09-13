// CreateOne 创建一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) CreateOne(ctx context.Context, data *{{.lowerDBName}}_model.{{.upperTableName}}) error {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}