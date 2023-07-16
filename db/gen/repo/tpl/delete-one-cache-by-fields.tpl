// DeleteOneCacheBy{{.upperFields}} 根据{{.lowerField}}删除一条数据并清理缓存
func (r *{{.upperTableName}}Repo) DeleteOneCacheBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) error {
	dao := {{.lowerDbName}}_dao.Use(r.db).{{.upperTableName}}
	first, err := dao.WithContext(ctx).Where({{.whereFields}}).First()
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if first == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where({{.whereFields}}).Delete()
	if err != nil {
		return err
	}
	err = r.DeleteUniqueIndexCache(ctx, []*{{.lowerDbName}}_model.{{.upperTableName}}{first})
	if err != nil {
		return err
	}
	return nil


}