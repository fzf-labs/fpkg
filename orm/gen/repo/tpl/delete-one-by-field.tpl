// DeleteOneBy{{.upperField}} 根据{{.lowerField}}删除一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) error {
	dao := {{.dbName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).Delete()
	if err != nil {
		return err
	}
	return nil
}