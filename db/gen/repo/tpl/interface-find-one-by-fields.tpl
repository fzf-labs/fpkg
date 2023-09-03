// FindOneBy{{.upperFields}} 根据{{.upperFields}}查询一条数据
FindOneBy{{.upperFields}}(ctx context.Context, {{.fieldAndDataTypes}}) (*{{.lowerDBName}}_model.{{.upperTableName}}, error)