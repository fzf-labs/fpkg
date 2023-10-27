// FindOneCacheBy{{.upperFields}} 根据{{.upperFields}}查询一条数据并设置缓存
FindOneCacheBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) (*{{.dbName}}_model.{{.upperTableName}}, error)