// CreateBatch 批量创建数据
CreateBatch(ctx context.Context, data []*{{.dbName}}_model.{{.upperTableName}}, batchSize int) error