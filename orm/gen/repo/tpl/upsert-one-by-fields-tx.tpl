// UpsertOneByFieldsTx Upsert一条数据，根据fields字段(事务)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpsertOneByFieldsTx(ctx context.Context,tx *{{.dbName}}_dao.Query, data *{{.dbName}}_model.{{.upperTableName}},fields []string) error {
	if len(fields) == 0 {
        return errors.New("UpsertOneByFieldsTx fields is empty")
    }
	columns := make([]clause.Column, 0)
	for _, v := range fields {
		columns = append(columns, clause.Column{Name: v})
	}
	dao := tx.{{.upperTableName}}
	err := dao.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   columns,
		UpdateAll: true,
	}).Create(data)
	if err != nil {
		return err
	}
	return nil
}