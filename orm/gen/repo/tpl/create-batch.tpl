// CreateBatch 批量创建数据
func ({{.firstTableChar}} *{{.upperTableName}}Repo) CreateBatch(ctx context.Context, data []*{{.lowerDBName}}_model.{{.upperTableName}}, batchSize int) error {
	dao := {{.lowerDBName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	err := dao.WithContext(ctx).CreateInBatches(data,batchSize)
	if err != nil {
		return err
	}
	return nil
}