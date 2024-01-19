// FindMultiByPaginator 查询分页数据(通用)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiByPaginator(ctx context.Context, paginatorReq *paginator.PaginatorReq) ([]*{{.dbName}}_model.{{.upperTableName}}, *paginator.PaginatorReply, error) {
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
	paginatorReply,err := paginatorReq.ConvertToPage(int32(total))
	if err != nil {
		return result, nil, err
	}
	query := {{.firstTableChar}}.db.WithContext(ctx).Model(&{{.dbName}}_model.{{.upperTableName}}{}).Clauses(whereExpressions...).Clauses(orderExpressions...)
	if paginatorReply.Page != 0 && paginatorReply.PageSize != 0 {
		query = query.Offset(int((paginatorReply.Page - 1) * paginatorReply.PageSize))
		query = query.Limit(int(paginatorReply.PageSize))
	}
	err = query.Find(&result).Error
	if err != nil {
		return result, nil, err
	}
	return result, paginatorReply, err
}
