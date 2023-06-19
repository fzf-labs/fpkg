package repo

import "gorm.io/gorm"

type MysqlModel struct {
	gorm   *gorm.DB
	schema string
}

func NewMysqlModel(gorm *gorm.DB) *MysqlModel {
	return &MysqlModel{gorm: gorm, schema: gorm.Migrator().CurrentDatabase()}
}

// FindAllTables 查询库中的所有表名称
func (m *MysqlModel) FindAllTables() ([]string, error) {
	var tables []string
	err := m.gorm.Debug().Raw("select table_name from information_schema.tables where table_schema= ?", m.schema).Scan(&tables).Error
	if err != nil {
		return nil, err
	}
	return tables, nil
}

// FindColumns 查询列
func (m *MysqlModel) FindColumns(table string) ([]*DbColumn, error) {
	sql := `SELECT ordinal_position as num,column_name as column_name,column_type AS data_type,character_set_name as character_set,collation_name as collation,is_nullable as is_nullable,column_default as column_default,extra as extra,column_name AS foreign_key,column_comment AS comment FROM information_schema.columns WHERE table_schema= ? AND table_name= ? ORDER BY ordinal_position asc`
	var reply []*DbColumn
	err := m.gorm.Raw(sql, m.schema, table).Scan(&reply).Error
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// FindIndex 查询索引
func (m *MysqlModel) FindIndex(table string) ([]*DbIndex, error) {
	querySql := `SELECT index_name as index_name,index_type AS index_algorithm,CASE non_unique WHEN 0 THEN'TRUE'ELSE'FALSE'END AS is_unique,column_name as column_name FROM information_schema.statistics WHERE table_schema= ? AND table_name= ? ORDER BY seq_in_index ASC;`
	reply := make([]*DbIndex, 0)
	err := m.gorm.Raw(querySql, m.schema, table).Scan(&reply).Error
	if err != nil {
		return nil, err
	}
	return reply, nil
}
