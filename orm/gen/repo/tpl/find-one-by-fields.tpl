// FindOneBy{{.upperFields}} 根据{{.upperFields}}查询一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindOneBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) (*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
    dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
    result, err := dao.WithContext(ctx).Where({{.whereFields}}).First()
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }
	return result, nil
}