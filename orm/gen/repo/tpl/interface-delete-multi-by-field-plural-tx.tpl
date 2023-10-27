// DeleteMultiBy{{.upperFieldPlural}} 根据{{.upperFieldPlural}}删除多条数据
DeleteMultiBy{{.upperFieldPlural}}Tx(ctx context.Context,tx *{{.dbName}}_dao.Query, {{.lowerFieldPlural}} []{{.dataType}}) error