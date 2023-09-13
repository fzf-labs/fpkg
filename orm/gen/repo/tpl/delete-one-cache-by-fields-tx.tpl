// DeleteOneCacheBy{{.upperFields}} 根据{{.lowerField}}删除一条数据并清理缓存
func ({{.firstTableChar}} *{{.upperTableName}}Repo) DeleteOneCacheBy{{.upperFields}}Tx(ctx context.Context,tx *{{.lowerDBName}}_dao.Query, {{.fieldAndDataTypes}}) error {
	dao := tx.{{.upperTableName}}
	first, err := dao.WithContext(ctx).Where({{.whereFields}}).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if first == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where({{.whereFields}}).Delete()
	if err != nil {
		return err
	}
	err = {{.firstTableChar}}.DeleteUniqueIndexCache(ctx, []*{{.lowerDBName}}_model.{{.upperTableName}}{first})
	if err != nil {
		return err
	}
	return nil
}