// FindOneBy{{.upperField}} 根据{{.lowerField}}查询一条数据并设置缓存
FindOneBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.lowerTableName}}_model.{{.upperTableName}}, error)