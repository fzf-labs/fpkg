package orm

import (
	"fmt"
	"math"
	"reflect"
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

// ConvertToGormExpression 根据SearchColumn参数转换为符合gorm where clause.Expression
func (p *PaginatorReq) ConvertToGormExpression(model any) (whereExpressions, orderExpressions []clause.Expression, err error) {
	whereExpressions = make([]clause.Expression, 0)
	orderExpressions = make([]clause.Expression, 0)
	jsonToColumn := p.jsonToColumn(model)
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if len(p.Search) > 0 {
		cols := make([]clause.Expression, 0)
		for _, v := range p.Search {
			if v.Field == "" {
				return whereExpressions, orderExpressions, fmt.Errorf("field cannot be empty")
			}
			if _, ok := jsonToColumn[v.Field]; !ok {
				return whereExpressions, orderExpressions, fmt.Errorf("field is not exist")
			}
			if v.Exp == "" {
				v.Exp = "="
			}
			if v.Logic == "" {
				v.Logic = "AND"
			}
			if _, ok := expMap[strings.TrimSpace(v.Exp)]; !ok {
				return whereExpressions, orderExpressions, fmt.Errorf("unknown s exp type '%s'", v.Exp)
			}
			if _, ok := logicMap[strings.TrimSpace(v.Logic)]; !ok {
				return whereExpressions, orderExpressions, fmt.Errorf("unknown logic type '%s'", v.Logic)
			}
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
					values := make([]any, 0)
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
				whereExpressions = append(whereExpressions, clause.And(cols...))
			} else {
				whereExpressions = append(whereExpressions, clause.Or(cols...))
			}
		}
	}
	if p.Order != "" {
		split := strings.Split(p.Order, ",")
		if len(split) != 2 {
			return whereExpressions, orderExpressions, fmt.Errorf("order format error")
		}
		if split[1] != "ASC" && split[1] != "DESC" {
			return whereExpressions, orderExpressions, fmt.Errorf("order format error")
		}
		desc := false
		if split[1] == "DESC" {
			desc = true
		}
		orderExpressions = append(orderExpressions, clause.OrderBy{
			Columns: []clause.OrderByColumn{
				{
					Column:  clause.Column{Name: split[0]},
					Desc:    desc,
					Reorder: false,
				},
			},
		})
	}
	return whereExpressions, orderExpressions, nil
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

// jsonToColumn 将model的tag中json和gorm的tag的Column转换为map[string]string
func (p *PaginatorReq) jsonToColumn(model any) map[string]string {
	m := make(map[string]string)
	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		json := field.Tag.Get("json")
		gorm := field.Tag.Get("gorm")
		if json != "" && gorm != "" {
			gorms := strings.Split(gorm, ";")
			for _, v := range gorms {
				if strings.Contains(v, "column") {
					column := strings.Split(v, ":")
					if len(column) == 2 {
						m[json] = column[1]
					}
				}
			}
		}
	}
	return m
}
