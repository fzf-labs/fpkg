// DeleteOneCacheBy{{.upperFields}} 根据{{.lowerField}}删除一条数据并清理缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteOneCacheBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) error {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	result, err := dao.WithContext(ctx).Where({{.whereFields}}).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if result == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where({{.whereFields}}).Delete()
	if err != nil {
		return err
	}
	err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, []*{{.lowerDBName}}_model.{{.upperTableName}}{result})
	if err != nil {
		return err
	}
	return nil
}