// FindOneBy{{.upperField}} 根据{{.lowerField}}查询一条数据
FindOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.lowerDbName}}_model.{{.upperTableName}}, error)