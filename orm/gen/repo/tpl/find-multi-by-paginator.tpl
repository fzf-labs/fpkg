// FindMultiByPaginator 查询分页数据(通用)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiByPaginator(ctx context.Context, paginatorReq *orm.PaginatorReq) ([]*{{.lowerDBName}}_model.{{.upperTableName}}, *orm.PaginatorReply, error) {
	result := make([]*{{.lowerDBName}}_model.{{.upperTableName}}, 0)
	var total int64
	whereExpressions, orderExpressions, err := paginatorReq.ConvertToGormExpression({{.lowerDBName}}_model.{{.upperTableName}}{})
	if err != nil {
		return nil, nil, err
	}
	err = {{.firstTableChar}}.db.WithContext(ctx).Model(&{{.lowerDBName}}_model.{{.upperTableName}}{}).Select([]string{"id"}).Clauses(whereExpressions...).Count(&total).Error
	if err != nil {
		return result, nil, err
	}
	if total == 0 {
		return result, nil, nil
	}
	paginatorReply := paginatorReq.ConvertToPage(int(total))
	err = {{.firstTableChar}}.db.WithContext(ctx).Model(&{{.lowerDBName}}_model.{{.upperTableName}}{}).Limit(paginatorReply.Limit).Offset(paginatorReply.Offset).Clauses(whereExpressions...).Clauses(orderExpressions...).Find(&result).Error
	if err != nil {
		return result, nil, err
	}
	return result, paginatorReply, err
}
