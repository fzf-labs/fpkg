package orm

import (
	"fmt"
	"math"
	"strings"
)

var expMap = map[string]string{
	"=":    " = ",
	"!=":   " <> ",
	">":    " > ",
	">=":   " >= ",
	"<":    " < ",
	"<=":   " <= ",
	"IN":   " IN ",
	"Like": " LIKE ",
}

var logicMap = map[string]string{
	"AND": " AND ",
	"OR":  " OR ",
}

type PaginatorReq struct {
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Order    string          `json:"order"`
	Search   []*SearchColumn `json:"search,omitempty"`
}

type PaginatorReply struct {
	Page      int `json:"Page"`
	PageSize  int `json:"PageSize"`
	Total     int `json:"Total"`
	PrevPage  int `json:"PrevPage"`
	NextPage  int `json:"NextPage"`
	TotalPage int `json:"TotalPage"`
	Limit     int `json:"-"`
	Offset    int `json:"-"`
}

type SearchColumn struct {
	Field string `json:"field"` // 字段
	Value string `json:"value"` // 值
	Exp   string `json:"exp"`   // 表达式 =, ! =, >, >=, <, <=, like
	Logic string `json:"logic"` // 逻辑关系 and,or
}

// convert 字段数据转换
func (s *SearchColumn) convert() error {
	if s.Field == "" {
		return fmt.Errorf("field 'name' cannot be empty")
	}
	if s.Exp == "" {
		s.Exp = "="
	}
	if v, ok := expMap[strings.TrimSpace(s.Exp)]; ok {
		s.Exp = v
		if s.Exp == " LIKE " {
			s.Value = fmt.Sprintf("%%%v%%", s.Value)
		}
	} else {
		return fmt.Errorf("unknown s expression type '%s'", s.Exp)
	}
	if s.Logic == "" {
		s.Logic = "AND"
	}
	if v, ok := logicMap[strings.TrimSpace(s.Logic)]; ok {
		s.Logic = v
	} else {
		return fmt.Errorf("unknown logic type '%s'", s.Logic)
	}
	return nil
}

// ConvertToGormConditions 根据SearchColumn参数转换为符合gorm的参数
func (p *PaginatorReq) ConvertToGormConditions() (str string, args []any, err error) {
	l := len(p.Search)
	if l == 0 {
		return "", nil, nil
	}
	for _, column := range p.Search {
		if column.Exp == "IN" {
			str = column.Field + " IN (?)" + column.Logic
			args = []any{args}
		} else {
			err := column.convert()
			if err != nil {
				return "", nil, err
			}
			str += column.Field + column.Exp + "?" + column.Logic
			args = append(args, column.Value)
		}
	}
	for _, v := range logicMap {
		str = strings.TrimRight(str, v)
	}
	return str, args, nil
}

// ConvertToPage 转换为page
func (p *PaginatorReq) ConvertToPage(total int) *PaginatorReply {
	// 根据nums总数，和prePage每页数量 生成分页总数
	page := p.Page
	pageSize := p.PageSize
	totalPage := int(math.Ceil(float64(total) / float64(p.PageSize))) // page总数
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
	return &PaginatorReply{
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

// ConvertToOrder 转换为page
func (p *PaginatorReq) ConvertToOrder() string {
	return p.Order
}
