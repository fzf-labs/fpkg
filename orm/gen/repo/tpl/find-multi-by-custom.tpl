// FindMultiByCustom 自定义查询数据(通用)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiByCustom(ctx context.Context, customReq *custom.PaginatorReq) ([]*{{.dbName}}_model.{{.upperTableName}}, *custom.PaginatorReply, error) {
	result := make([]*{{.dbName}}_model.{{.upperTableName}}, 0)
	var total int64
	whereExpressions, orderExpressions, err := customReq.ConvertToGormExpression({{.dbName}}_model.{{.upperTableName}}{})
	if err != nil {
		return result, nil, err
	}
	err = {{.firstTableChar}}.db.WithContext(ctx).Model(&{{.dbName}}_model.{{.upperTableName}}{}).Select([]string{"*"}).Clauses(whereExpressions...).Count(&total).Error
	if err != nil {
		return result, nil, err
	}
	if total == 0 {
		return result, nil, nil
	}
	customReply,err := customReq.ConvertToPage(int32(total))
	if err != nil {
		return result, nil, err
	}
	query := {{.firstTableChar}}.db.WithContext(ctx).Model(&{{.dbName}}_model.{{.upperTableName}}{}).Clauses(whereExpressions...).Clauses(orderExpressions...)
	if customReply.Page != 0 && customReply.PageSize != 0 {
		query = query.Offset(int((customReply.Page - 1) * customReply.PageSize))
		query = query.Limit(int(customReply.PageSize))
	}
	err = query.Find(&result).Error
	if err != nil {
		return result, nil, err
	}
	return result, customReply, err
}
