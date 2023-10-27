// FindOneBy{{.upperField}} 根据{{.lowerField}}查询一条数据
FindOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.dbName}}_model.{{.upperTableName}}, error)