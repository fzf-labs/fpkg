// FindOneBy{{.upperFields}} 根据{{.upperFields}}查询一条数据
func (r *{{.upperTableName}}Repo) FindOneBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) (*{{.lowerDbName}}_model.{{.upperTableName}}, error) {
    dao := {{.lowerDbName}}_dao.Use(r.db).{{.upperTableName}}
    result, err := dao.WithContext(ctx).Where({{.whereFields}}).First()
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }
	return result, nil
}