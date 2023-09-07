// DeleteOneBy{{.upperFields}} 根据{{.lowerField}}删除一条数据
func (r *{{.upperTableName}}Repo) DeleteOneBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) error {
	dao := {{.lowerDBName}}_dao.Use(r.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where({{.whereFields}}).Delete()
	if err != nil {
		return err
	}
	return nil
}