// CreateOneByTx 创建一条数据(事务)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) CreateOneByTx(ctx context.Context, tx *{{.lowerDBName}}_dao.Query, data *{{.lowerDBName}}_model.{{.upperTableName}}) error {
	dao := tx.{{.upperTableName}}
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}