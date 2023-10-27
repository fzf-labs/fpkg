// CreateBatch 批量创建数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) CreateBatch(ctx context.Context, data []*{{.dbName}}_model.{{.upperTableName}}, batchSize int) error {
	dao := {{.dbName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	err := dao.WithContext(ctx).CreateInBatches(data,batchSize)
	if err != nil {
		return err
	}
	return nil
}