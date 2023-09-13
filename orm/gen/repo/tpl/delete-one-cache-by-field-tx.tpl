// DeleteOneCacheBy{{.upperField}} 根据{{.lowerField}}删除一条数据并清理缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteOneCacheBy{{.upperField}}Tx(ctx context.Context,tx *{{.lowerDBName}}_dao.Query, {{.lowerField}} {{.dataType}}) error {
	dao := tx.{{.upperTableName}}
	first, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if first == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).Delete()
	if err != nil {
		return err
	}
	err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, []*{{.lowerDBName}}_model.{{.upperTableName}}{first})
    if err != nil {
    	return err
    }
	return nil
}