// FindOneBy{{.upperField}} 根据{{.lowerField}}查询一条数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.lowerDBName}}_model.{{.upperTableName}}, error) {
    dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
    result, err := dao.WithContext(ctx).Where(dao.{{.upperField}}.Eq({{.lowerField}})).First()
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, err
    }
	return result, nil
}