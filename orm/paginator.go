package orm

import (
	"fmt"
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

type PaginatorParams struct {
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Order    string          `json:"order"`
	Search   []*SearchColumn `json:"search,omitempty"`
}

// ConvertToGormConditions 根据SearchColumn参数转换为符合gorm的参数
func (p *PaginatorParams) ConvertToGormConditions() (string, []interface{}, error) {
	str := ""
	var args []interface{}
	l := len(p.Search)
	if l == 0 {
		return "", nil, nil
	}
	for _, column := range p.Search {
		if column.Exp == "IN" {
			str = column.Field + " IN (?)" + column.Logic
			args = []interface{}{args}
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
func (p *PaginatorParams) ConvertToPage() (limit int, offset int) {
	return p.PageSize, (p.Page - 1) * p.PageSize
}

// ConvertToOrder 转换为page
func (p *PaginatorParams) ConvertToOrder() string {
	if p.Order == "" {
		return "id DESC"
	}
	return p.Order
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
	if v, ok := expMap[strings.ToLower(s.Exp)]; ok {
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
	if v, ok := logicMap[strings.ToLower(s.Logic)]; ok {
		s.Logic = v
	} else {
		return fmt.Errorf("unknown logic type '%s'", s.Logic)
	}
	return nil
}
