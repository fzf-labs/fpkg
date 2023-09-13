// DeleteOneCacheBy{{.upperField}} 根据{{.lowerField}}删除一条数据并清理缓存
DeleteOneCacheBy{{.upperField}}Tx(ctx context.Context,tx *{{.lowerDBName}}_dao.Query, {{.lowerField}} {{.dataType}}) error