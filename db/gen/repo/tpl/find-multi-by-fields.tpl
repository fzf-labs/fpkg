// FindMultiBy{{.upperFields}} 根据{{.upperFields}}查询多条数据
func (r *{{.upperTableName}}Repo) FindMultiBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) ([]*{{.lowerDbName}}_model.{{.upperTableName}}, error) {
	dao := {{.lowerDbName}}_dao.Use(r.db).{{.upperTableName}}
	result, err := dao.WithContext(ctx).Where({{.whereFields}}).Find()
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }
	return result, nil
}
