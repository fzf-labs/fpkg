// DeleteOneBy{{.upperFields}} 根据{{.lowerField}}删除一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteOneBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) error {
	dao := {{.dbName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where({{.whereFields}}).Delete()
	if err != nil {
		return err
	}
	return nil
}