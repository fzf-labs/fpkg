// DeleteOneCacheBy{{.upperField}} 根据{{.lowerField}}删除一条数据并清理缓存
func (r *{{.lowerTableName}}Repo) DeleteOneCacheBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) error {
	dao := {{.lowerTableName}}_dao.Use(r.db).{{.upperTableName}}
	first, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).First()
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if first == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).Delete()
	if err != nil {
		return err
	}
	err = r.DeleteUniqueIndexCache(ctx, []*{{.lowerTableName}}_model.{{.upperTableName}}{first})
    if err != nil {
    	return err
    }
	return nil
}