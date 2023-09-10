// DeleteOneBy{{.upperField}} 根据{{.lowerField}}删除一条数据
func (r *{{.upperTableName}}Repo) DeleteOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) error {
	dao := {{.lowerDBName}}_dao.Use(r.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).Delete()
	if err != nil {
		return err
	}
	return nil
}