package paginator

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"gorm.io/gorm/clause"
)

type EXP string

const (
	Eq      EXP = "="
	Neq     EXP = "!="
	Gt      EXP = ">"
	Gte     EXP = ">="
	Lt      EXP = "<"
	Lte     EXP = "<="
	In      EXP = "IN"
	NotIn   EXP = "NotIn"
	Like    EXP = "Like"
	NotLike EXP = "NotLike"
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

type Req struct {
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Order    []*OrderColumn  `json:"order"`
	Search   []*SearchColumn `json:"search"`
}

type Reply struct {
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
	Exp   ORDER  `json:"exp"`   // 表达式 ASC,DESC
}

type SearchColumn struct {
	Field string `json:"field"` // 字段
	Value string `json:"value"` // 值
	Exp   EXP    `json:"exp"`   // 表达式 exp
	Logic LOGIC  `json:"logic"` // 逻辑关系 logic
}

// ConvertToGormExpression 根据SearchColumn参数转换为符合gorm where clause.Expression
func (p *Req) ConvertToGormExpression(model any) (whereExpressions, orderExpressions []clause.Expression, err error) {
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
			if !v.Exp.Validate() {
				return whereExpressions, orderExpressions, fmt.Errorf("unknown s exp type '%s'", v.Exp)
			}
			if v.Logic == "" {
				v.Logic = And
			}
			if !v.Logic.Validate() {
				return whereExpressions, orderExpressions, fmt.Errorf("unknown s logic type '%s'", v.Logic)
			}
			var expression clause.Expression
			switch v.Exp {
			case Eq:
				expression = clause.Eq{Column: jsonToColumn[v.Field], Value: v.Value}
			case Neq:
				expression = clause.Neq{Column: jsonToColumn[v.Field], Value: v.Value}
			case Gt:
				expression = clause.Gt{Column: jsonToColumn[v.Field], Value: v.Value}
			case Gte:
				expression = clause.Gte{Column: jsonToColumn[v.Field], Value: v.Value}
			case Lt:
				expression = clause.Lt{Column: jsonToColumn[v.Field], Value: v.Value}
			case Lte:
				expression = clause.Lte{Column: jsonToColumn[v.Field], Value: v.Value}
			case In:
				split := strings.Split(v.Value, ",")
				if len(split) > 0 {
					values := make([]any, 0)
					for _, vv := range split {
						values = append(values, vv)
					}
					expression = clause.IN{Column: jsonToColumn[v.Field], Values: values}
				}
			case NotIn:
				split := strings.Split(v.Value, ",")
				if len(split) > 0 {
					values := make([]any, 0)
					for _, vv := range split {
						values = append(values, vv)
					}
					expression = clause.Not(clause.IN{Column: jsonToColumn[v.Field], Values: values})
				}
			case Like:
				expression = clause.Like{Column: jsonToColumn[v.Field], Value: v.Value}
			case NotLike:
				expression = clause.Not(clause.Like{Column: jsonToColumn[v.Field], Value: v.Value})
			default:
				return whereExpressions, orderExpressions, fmt.Errorf("unknown s exp type '%s'", v.Exp)
			}
			if v.Logic == And {
				whereExpressions = append(whereExpressions, expression)
			} else {
				whereExpressions = append(whereExpressions, expression)
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
			if v.Exp == "" {
				v.Exp = Asc
			}
			if !v.Exp.Validate() {
				return whereExpressions, orderExpressions, fmt.Errorf("order is err")
			}
			orderExpressions = append(orderExpressions, clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{
						Column:  clause.Column{Name: jsonToColumn[v.Field]},
						Desc:    v.Exp == Desc,
						Reorder: false,
					},
				},
			})
		}
	}
	return whereExpressions, orderExpressions, nil
}

// ConvertToPage 转换为page
func (p *Req) ConvertToPage(total int) (*Reply, error) {
	resp := &Reply{
		Total: total,
	}
	if p.Page < 0 {
		return resp, fmt.Errorf("page cannot be less than 0")
	}
	if p.PageSize < 0 {
		return resp, fmt.Errorf("pageSize cannot be less than 0")
	}
	if (p.Page != 0 && p.PageSize == 0) || (p.Page == 0 && p.PageSize != 0) {
		return resp, fmt.Errorf("page and pageSize must be a pair")
	}
	if p.Page == 0 && p.PageSize == 0 {
		return resp, nil
	}
	resp.Page = p.Page
	resp.PageSize = p.PageSize
	resp.TotalPage = int(math.Ceil(float64(total) / float64(p.PageSize)))
	resp.NextPage = p.Page + 1
	if resp.NextPage > resp.TotalPage {
		resp.NextPage = resp.TotalPage
	}
	resp.PrevPage = p.Page - 1
	if resp.PrevPage <= 0 {
		resp.PrevPage = 1
	}
	resp.Limit = p.PageSize
	resp.Offset = (p.Page - 1) * p.PageSize
	return resp, nil
}

// jsonToColumn 将model的tag中json和gorm的tag的Column转换为map[string]string
func (p *Req) jsonToColumn(model any) map[string]string {
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
