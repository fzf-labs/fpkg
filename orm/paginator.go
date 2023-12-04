package orm

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"gorm.io/gorm/clause"
)

type EXP string

const (
	Eq   EXP = "="
	Neq  EXP = "!="
	Gt   EXP = ">"
	Gte  EXP = ">="
	Lt   EXP = "<"
	Lte  EXP = "<="
	In   EXP = "IN"
	Like EXP = "Like"
)

func (s EXP) Validate() bool {
	switch s {
	case Eq, Neq, Gt, Gte, Lt, Lte, In, Like:
		return true
	default:
		return false
	}
}

type LOGIC string

const (
	And LOGIC = "AND"
	Or  LOGIC = "OR"
)

func (s LOGIC) Validate() bool {
	switch s {
	case And, Or:
		return true
	default:
		return false
	}
}

type ORDER string

const (
	Asc  ORDER = "ASC"
	Desc ORDER = "DESC"
)

func (s ORDER) Validate() bool {
	switch s {
	case Asc, Desc:
		return true
	default:
		return false
	}
}

type PaginatorReq struct {
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Order    []*OrderColumn  `json:"order"`
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

type OrderColumn struct {
	Field string `json:"field"` // 字段
	Order ORDER  `json:"exp"`   // 表达式 ASC,DESC
}

type SearchColumn struct {
	Field string `json:"field"` // 字段
	Value string `json:"value"` // 值
	Exp   EXP    `json:"exp"`   // 表达式 exp
	Logic LOGIC  `json:"logic"` // 逻辑关系 logic
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
				v.Exp = Eq
			}
			if v.Logic == "" {
				v.Logic = And
			}
			if !v.Exp.Validate() {
				return whereExpressions, orderExpressions, fmt.Errorf("exp is err")
			}
			if !v.Logic.Validate() {
				return whereExpressions, orderExpressions, fmt.Errorf("logic is err")
			}
			if v.Exp == Eq {
				cols = append(cols, clause.Eq{Column: jsonToColumn[v.Field], Value: v.Value})
			}
			if v.Exp == Neq {
				cols = append(cols, clause.Neq{Column: jsonToColumn[v.Field], Value: v.Value})
			}
			if v.Exp == Gt {
				cols = append(cols, clause.Gt{Column: jsonToColumn[v.Field], Value: v.Value})
			}
			if v.Exp == Gte {
				cols = append(cols, clause.Gte{Column: jsonToColumn[v.Field], Value: v.Value})
			}
			if v.Exp == Lt {
				cols = append(cols, clause.Lt{Column: jsonToColumn[v.Field], Value: v.Value})
			}
			if v.Exp == Lte {
				cols = append(cols, clause.Lte{Column: jsonToColumn[v.Field], Value: v.Value})
			}
			if v.Exp == In {
				split := strings.Split(v.Value, ",")
				if len(split) > 0 {
					values := make([]any, 0)
					for _, vv := range split {
						values = append(values, vv)
					}
					cols = append(cols, clause.IN{Column: jsonToColumn[v.Field], Values: values})
				}
			}
			if v.Exp == Like {
				cols = append(cols, clause.Like{Column: jsonToColumn[v.Field], Value: v.Value})
			}
			if v.Logic == And {
				whereExpressions = append(whereExpressions, clause.And(cols...))
			} else {
				whereExpressions = append(whereExpressions, clause.Or(cols...))
			}
		}
	}
	if len(p.Order) > 0 {
		for _, v := range p.Order {
			if v.Field == "" {
				return whereExpressions, orderExpressions, fmt.Errorf("field cannot be empty")
			}
			if _, ok := jsonToColumn[v.Field]; !ok {
				return whereExpressions, orderExpressions, fmt.Errorf("field is not exist")
			}
			if v.Order == "" {
				v.Order = Asc
			}
			if !v.Order.Validate() {
				return whereExpressions, orderExpressions, fmt.Errorf("order is err")
			}
			orderExpressions = append(orderExpressions, clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{
						Column:  clause.Column{Name: jsonToColumn[v.Field]},
						Desc:    v.Order == Desc,
						Reorder: false,
					},
				},
			})
		}
	}
	return whereExpressions, orderExpressions, nil
}

// ConvertToPage 转换为page
func (p *PaginatorReq) ConvertToPage(total int) *PaginatorReply {
	// 根据nums总数，和prePage每页数量 生成分页总数
	page := p.Page
	pageSize := p.PageSize
	totalPage := int(math.Ceil(float64(total) / float64(p.PageSize))) // page总数
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
