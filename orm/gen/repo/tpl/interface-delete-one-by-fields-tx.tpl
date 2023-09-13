// DeleteOneBy{{.upperFields}} 根据{{.upperFields}}删除一条数据
DeleteOneBy{{.upperFields}}Tx(ctx context.Context,tx *{{.lowerDBName}}_dao.Query, {{.fieldAndDataTypes}}) error