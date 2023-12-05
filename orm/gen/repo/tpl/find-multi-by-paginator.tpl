// FindMultiByPaginator 查询分页数据(通用)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiByPaginator(ctx context.Context, paginatorReq *paginator.Req) ([]*{{.dbName}}_model.{{.upperTableName}}, *paginator.Reply, error) {
	result := make([]*{{.dbName}}_model.{{.upperTableName}}, 0)
	var total int64
	whereExpressions, orderExpressions, err := paginatorReq.ConvertToGormExpression({{.dbName}}_model.{{.upperTableName}}{})
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
	paginatorReply,err := paginatorReq.ConvertToPage(int(total))
	if err != nil {
		return result, nil, err
	}
	query := {{.firstTableChar}}.db.WithContext(ctx).Model(&{{.dbName}}_model.{{.upperTableName}}{}).Clauses(whereExpressions...).Clauses(orderExpressions...)
	if paginatorReply.Offset != 0 {
		query = query.Offset(paginatorReply.Offset)
	}
	if paginatorReply.Limit != 0 {
		query = query.Limit(paginatorReply.Limit)
	}
	err = query.Find(&result).Error
	if err != nil {
		return result, nil, err
	}
	return result, paginatorReply, err
}
