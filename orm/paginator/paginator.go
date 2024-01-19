package paginator

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"gorm.io/gorm/clause"
)

const (
	EQ      string = "="
	NEQ     string = "!="
	GT      string = ">"
	GTE     string = ">="
	LT      string = "<"
	LTE     string = "<="
	IN      string = "IN"
	NOTIN   string = "NOT IN"
	LIKE    string = "LIKE"
	NOTLIKE string = "NOT LIKE"
)

const (
	AND string = "AND"
	OR  string = "OR"
)

const (
	ASC  string = "ASC"
	DESC string = "DESC"
)

func ExpValidate(s string) bool {
	switch s {
	case EQ, NEQ, GT, GTE, LT, LTE, IN, LIKE:
		return true
	default:
		return false
	}
}

func LogicValidate(s string) bool {
	switch s {
	case AND, OR:
		return true
	default:
		return false
	}
}

func OrderValidate(s string) bool {
	switch s {
	case ASC, DESC:
		return true
	default:
		return false
	}
}

// ConvertToGormExpression 根据SearchColumn参数转换为符合gorm where clause.Expression
func (p *PaginatorReq) ConvertToGormExpression(model any) (whereExpressions, orderExpressions []clause.Expression, err error) {
	whereExpressions = make([]clause.Expression, 0)
	orderExpressions = make([]clause.Expression, 0)
	column := jsonToColumn(model)
	if len(p.Search) > 0 {
		for _, v := range p.Search {
			if v.Field == "" {
				return whereExpressions, orderExpressions, fmt.Errorf("field cannot be empty")
			}
			if _, ok := column[v.Field]; !ok {
				return whereExpressions, orderExpressions, fmt.Errorf("field is not exist")
			}
			if v.Exp == "" {
				v.Exp = EQ
			}
			if !ExpValidate(v.Exp) {
				return whereExpressions, orderExpressions, fmt.Errorf("unknown s exp type '%s'", v.Exp)
			}
			if v.Logic == "" {
				v.Logic = AND
			}
			if !LogicValidate(v.Logic) {
				return whereExpressions, orderExpressions, fmt.Errorf("unknown s logic type '%s'", v.Logic)
			}
			var expression clause.Expression
			switch v.Exp {
			case EQ:
				expression = clause.Eq{Column: column[v.Field], Value: v.Value}
			case NEQ:
				expression = clause.Neq{Column: column[v.Field], Value: v.Value}
			case GT:
				expression = clause.Gt{Column: column[v.Field], Value: v.Value}
			case GTE:
				expression = clause.Gte{Column: column[v.Field], Value: v.Value}
			case LT:
				expression = clause.Lt{Column: column[v.Field], Value: v.Value}
			case LTE:
				expression = clause.Lte{Column: column[v.Field], Value: v.Value}
			case IN:
				split := strings.Split(v.Value, ",")
				if len(split) > 0 {
					values := make([]any, 0)
					for _, vv := range split {
						values = append(values, vv)
					}
					expression = clause.IN{Column: column[v.Field], Values: values}
				}
			case NOTIN:
				split := strings.Split(v.Value, ",")
				if len(split) > 0 {
					values := make([]any, 0)
					for _, vv := range split {
						values = append(values, vv)
					}
					expression = clause.Not(clause.IN{Column: column[v.Field], Values: values})
				}
			case LIKE:
				expression = clause.Like{Column: column[v.Field], Value: v.Value}
			case NOTLIKE:
				expression = clause.Not(clause.Like{Column: column[v.Field], Value: v.Value})
			default:
				return whereExpressions, orderExpressions, fmt.Errorf("unknown s exp type '%s'", v.Exp)
			}
			if v.Logic == AND {
				whereExpressions = append(whereExpressions, clause.And(expression))
			} else {
				whereExpressions = append(whereExpressions, clause.Or(expression))
			}
		}
	}
	if len(p.Order) > 0 {
		for _, v := range p.Order {
			if v.Field == "" {
				return whereExpressions, orderExpressions, fmt.Errorf("field cannot be empty")
			}
			if _, ok := column[v.Field]; !ok {
				return whereExpressions, orderExpressions, fmt.Errorf("field is not exist")
			}
			if v.Exp == "" {
				v.Exp = ASC
			}
			if !OrderValidate(v.Exp) {
				return whereExpressions, orderExpressions, fmt.Errorf("order is err")
			}
			orderExpressions = append(orderExpressions, clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{
						Column:  clause.Column{Name: column[v.Field]},
						Desc:    v.Exp == DESC,
						Reorder: false,
					},
				},
			})
		}
	}
	return whereExpressions, orderExpressions, nil
}

// ConvertToPage 转换为page
func (p *PaginatorReq) ConvertToPage(total int32) (*PaginatorReply, error) {
	resp := &PaginatorReply{
		Page:      0,
		PageSize:  0,
		Total:     total,
		PrevPage:  0,
		NextPage:  0,
		TotalPage: 0,
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
	resp.TotalPage = int32(math.Ceil(float64(total) / float64(p.PageSize)))
	resp.NextPage = p.Page + 1
	if resp.NextPage > resp.TotalPage {
		resp.NextPage = resp.TotalPage
	}
	resp.PrevPage = p.Page - 1
	if resp.PrevPage <= 0 {
		resp.PrevPage = 1
	}
	return resp, nil
}

// jsonToColumn 将model的tag中json和gorm的tag的Column转换为map[string]string
func jsonToColumn(model any) map[string]string {
	m := make(map[string]string)
	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		json := field.Tag.Get("json")
		gorm := field.Tag.Get("gorm")
		if json != "" && gorm != "" {
			gormTags := strings.Split(gorm, ";")
			for _, v := range gormTags {
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
