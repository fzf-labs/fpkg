// DeleteOneCacheBy{{.upperFields}} 根据{{.upperFields}}删除一条数据并清理缓存
DeleteOneCacheBy{{.upperFields}}Tx(ctx context.Context,tx *{{.lowerDBName}}_dao.Query, {{.fieldAndDataTypes}}) error