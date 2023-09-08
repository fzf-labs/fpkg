// CreateBatch 批量创建数据
CreateBatch(ctx context.Context, data []*{{.lowerDBName}}_model.{{.upperTableName}}, batchSize int) error