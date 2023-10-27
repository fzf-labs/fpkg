// DeleteOneBy{{.upperField}} 根据{{.lowerField}}删除一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteOneBy{{.upperField}}Tx(ctx context.Context,tx *{{.dbName}}_dao.Query, {{.lowerField}} {{.dataType}}) error {
	dao := tx.{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).Delete()
	if err != nil {
		return err
	}
	return nil
}