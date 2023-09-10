// DeleteMultiBy{{.upperFieldPlural}} 根据{{.lowerFieldPlural}}删除多条数据
func (r *{{.upperTableName}}Repo) DeleteMultiBy{{.upperFieldPlural}}(ctx context.Context, {{.lowerFieldPlural}} []{{.dataType}}) error {
	dao := {{.lowerDBName}}_dao.Use(r.db).{{.upperTableName}}
	_, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.In({{.lowerFieldPlural}}...)).Delete()
	if err != nil {
		return err
	}
	return nil
}
