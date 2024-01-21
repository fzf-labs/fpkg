// FindMultiByCustom 自定义查询数据(通用)
FindMultiByCustom(ctx context.Context, customReq *custom.PaginatorReq) ([]*{{.dbName}}_model.{{.upperTableName}}, *custom.PaginatorReply, error)