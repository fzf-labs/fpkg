// SaveOne 保存一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) SaveOne(ctx context.Context, data *{{.lowerDBName}}_model.{{.upperTableName}}) error {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	err := dao.WithContext(ctx).Save(data)
	if err != nil {
		return err
	}
	return nil
}