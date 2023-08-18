package page

import "math"

type Page struct {
	Page      int `json:"Page"`
	PageSize  int `json:"PageSize"`
	Total     int `json:"Total"`
	PrevPage  int `json:"PrevPage"`
	NextPage  int `json:"NextPage"`
	TotalPage int `json:"TotalPage"`
	Limit     int `json:"-"`
	Offset    int `json:"-"`
}

func Paginator(page, pageSize, total int) *Page {
	// 根据nums总数，和prePage每页数量 生成分页总数
	totalPage := int(math.Ceil(float64(total) / float64(pageSize))) // page总数
	if page > totalPage {
		page = totalPage
	}
	if page <= 0 {
		page = 1
	}
	prevPage := page - 1
	if prevPage <= 0 {
		prevPage = 1
	}
	nextPage := page + 1
	if nextPage > totalPage {
		nextPage = totalPage
	}
	return &Page{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		PrevPage:  prevPage,
		NextPage:  nextPage,
		TotalPage: totalPage,
		Limit:     pageSize,
		Offset:    (page - 1) * pageSize,
	}
}
