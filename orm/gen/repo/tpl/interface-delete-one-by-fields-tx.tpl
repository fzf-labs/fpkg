// DeleteOneBy{{.upperFields}} 根据{{.upperFields}}删除一条数据
DeleteOneBy{{.upperFields}}Tx(ctx context.Context,tx *{{.dbName}}_dao.Query, {{.fieldAndDataTypes}}) error