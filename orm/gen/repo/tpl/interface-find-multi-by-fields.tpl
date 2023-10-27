// FindMultiBy{{.upperFields}} 根据{{.upperFields}}查询多条数据
FindMultiBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) ([]*{{.dbName}}_model.{{.upperTableName}}, error)