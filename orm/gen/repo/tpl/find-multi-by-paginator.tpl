// FindMultiByPaginator 查询分页数据(通用)
func ({{.firstTableChar}} *{{.upperTableName}}Repo) FindMultiByPaginator(ctx context.Context, params *orm.PaginatorParams) ([]*{{.lowerDBName}}_model.{{.upperTableName}}, int64, error) {
	result := make([]*{{.lowerDBName}}_model.{{.upperTableName}}, 0)
	var total int64
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, err
	}
	err = {{.firstTableChar}}.db.WithContext(ctx).Model(&{{.lowerDBName}}_model.{{.upperTableName}}{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, total, nil
	}
	query := {{.firstTableChar}}.db.WithContext(ctx)
	order := params.ConvertToOrder()
	if(order != ""){
	    query = query.Order(order)
	}
	limit, offset := params.ConvertToPage()
	err = query.Limit(limit).Offset(offset).Where(queryStr, args...).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, total, err
}
