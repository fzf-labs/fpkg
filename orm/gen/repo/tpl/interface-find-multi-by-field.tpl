// FindMultiBy{{.upperField}} 根据{{.lowerField}}查询多条数据
FindMultiBy{{.upperField}}(ctx context.Context, {{.lowerField}} {{.dataType}}) ([]*{{.dbName}}_model.{{.upperTableName}}, error)