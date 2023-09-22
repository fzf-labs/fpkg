// FindMultiByPaginator 查询分页数据(通用)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiByPaginator(ctx context.Context, paginatorReq *orm.PaginatorReq) ([]*{{.lowerDBName}}_model.{{.upperTableName}}, *orm.PaginatorReply, error) {
	result := make([]*{{.lowerDBName}}_model.{{.upperTableName}}, 0)
	var total int64
	queryStr, args, err := paginatorReq.ConvertToGormConditions()
	if err != nil {
		return nil, nil, err
	}
	err = {{.firstTableChar}}.db.WithContext(ctx).Model(&{{.lowerDBName}}_model.{{.upperTableName}}{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
	if err != nil {
		return nil, nil, err
	}
	if total == 0 {
		return nil, nil, nil
	}
	query := {{.firstTableChar}}.db.WithContext(ctx)
	order := paginatorReq.ConvertToOrder()
	if(order != ""){
	    query = query.Order(order)
	}
	paginatorReply := paginatorReq.ConvertToPage(int(total))
	err = query.Limit(paginatorReply.Limit).Offset(paginatorReply.Offset).Where(queryStr, args...).Find(&result).Error
	if err != nil {
		return nil, nil, err
	}
	return result, paginatorReply, err
}
