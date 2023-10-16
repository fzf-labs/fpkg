package orm

import (
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm/clause"
)

var expMap = map[string]struct{}{
	"=":    {},
	"!=":   {},
	">":    {},
	">=":   {},
	"<":    {},
	"<=":   {},
	"IN":   {},
	"Like": {},
}

var logicMap = map[string]struct{}{
	"AND": {},
	"OR":  {},
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
	Exp   string `json:"exp"`   // 表达式 expMap
	Logic string `json:"logic"` // 逻辑关系 logicMap
}

// Check 字段校验
func (p *PaginatorReq) Check() error {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.Order != "" {
		split := strings.Split(p.Order, ",")
		if len(split) != 2 {
			return fmt.Errorf("order format error")
		}
		if split[1] != "ASC" && split[1] != "DESC" {
			return fmt.Errorf("order format error")
		}
	}
	if len(p.Search) > 0 {
		for k := range p.Search {
			if p.Search[k].Field == "" {
				return fmt.Errorf("field 'name' cannot be empty")
			}
			if p.Search[k].Exp == "" {
				p.Search[k].Exp = "="
			}
			if p.Search[k].Logic == "" {
				p.Search[k].Logic = "AND"
			}
			if _, ok := expMap[strings.TrimSpace(p.Search[k].Exp)]; !ok {
				return fmt.Errorf("unknown s exp type '%s'", p.Search[k].Exp)
			}
			if _, ok := logicMap[strings.TrimSpace(p.Search[k].Logic)]; !ok {
				return fmt.Errorf("unknown logic type '%s'", p.Search[k].Logic)
			}
		}
	}
	return nil
}

// ConvertToGormWhereExpression 根据SearchColumn参数转换为符合gorm where clause.Expression
func (p *PaginatorReq) ConvertToGormWhereExpression() []clause.Expression {
	expressions := make([]clause.Expression, 0)
	l := len(p.Search)
	if l == 0 {
		return expressions
	}
	cols := make([]clause.Expression, 0)
	for _, v := range p.Search {
		if v.Exp == "=" {
			cols = append(cols, clause.Eq{Column: v.Field, Value: v.Value})
		}
		if v.Exp == "!=" {
			cols = append(cols, clause.Neq{Column: v.Field, Value: v.Value})
		}
		if v.Exp == ">" {
			cols = append(cols, clause.Gt{Column: v.Field, Value: v.Value})
		}
		if v.Exp == ">=" {
			cols = append(cols, clause.Gte{Column: v.Field, Value: v.Value})
		}
		if v.Exp == "<" {
			cols = append(cols, clause.Lt{Column: v.Field, Value: v.Value})
		}
		if v.Exp == "<=" {
			cols = append(cols, clause.Lte{Column: v.Field, Value: v.Value})
		}
		if v.Exp == "IN" {
			split := strings.Split(v.Value, ",")
			if len(split) > 0 {
				values := make([]interface{}, 0)
				for _, vv := range split {
					values = append(values, vv)
				}
				cols = append(cols, clause.IN{Column: v.Field, Values: values})
			}
		}
		if v.Exp == "Like" {
			cols = append(cols, clause.Like{Column: v.Field, Value: v.Value})
		}
		if v.Logic == "AND" {
			expressions = append(expressions, clause.And(cols...))
		} else {
			expressions = append(expressions, clause.Or(cols...))
		}
	}
	return expressions
}

// ConvertToGormOrderExpression 根据SearchColumn参数转换为符合gorm order clause.Expression
func (p *PaginatorReq) ConvertToGormOrderExpression() []clause.Expression {
	expressions := make([]clause.Expression, 0)
	if p.Order != "" {
		split := strings.Split(p.Order, ",")
		desc := false
		if split[1] == "DESC" {
			desc = true
		}
		expressions = append(expressions, clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column:  clause.Column{Name: split[0]},
					Desc:    desc,
					Reorder: false,
				},
			},
		})
	}
	return expressions
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
