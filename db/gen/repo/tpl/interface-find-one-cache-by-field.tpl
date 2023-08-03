// FindOneCacheBy{{.upperField}} 根据{{.lowerField}}查询一条数据并设置缓存
FindOneCacheBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) (*{{.lowerDbName}}_model.{{.upperTableName}}, error)