// FindMultiBy{{.upperFields}} 根据{{.upperFields}}查询多条数据
FindMultiBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) ([]*{{.lowerDBName}}_model.{{.upperTableName}}, error)