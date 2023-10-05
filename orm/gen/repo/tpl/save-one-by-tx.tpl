// SaveOneByTx 保存一条数据(事务)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) SaveOneByTx(ctx context.Context, tx *{{.lowerDBName}}_dao.Query, data *{{.lowerDBName}}_model.{{.upperTableName}}) error {
	dao := tx.{{.upperTableName}}
	err := dao.WithContext(ctx).Save(data)
	if err != nil {
		return err
	}
	return nil
}