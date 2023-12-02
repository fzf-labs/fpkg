// UpsertOneByFields Upsert一条数据，根据fields字段
func ({{.firstTableChar}} *{{.upperTableName}}Repo) UpsertOneByFields(ctx context.Context, data *{{.dbName}}_model.{{.upperTableName}},fields []string) error {
	if len(fields) == 0 {
        return errors.New("UpsertOneByFields fields is empty")
    }
	columns := make([]clause.Column, 0)
	for _, v := range fields {
		columns = append(columns, clause.Column{Name: v})
	}
	dao := {{.dbName}}_dao.Use({{.firstTableChar}}.db).{{.upperTableName}}
	err := dao.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   columns,
		UpdateAll: true,
	}).Create(data)
	if err != nil {
		return err
	}
	return nil
}