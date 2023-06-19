package repo

import (
	"gorm.io/gorm"
)

type PostgresqlModel struct {
	gorm   *gorm.DB
	schema string
}

func NewPostgresqlModel(gorm *gorm.DB) *PostgresqlModel {
	return &PostgresqlModel{gorm: gorm, schema: "public"}
}

// FindAllTables 查询库中的所有表名称
func (m *PostgresqlModel) FindAllTables() ([]string, error) {
	var tables []string
	err := m.gorm.Debug().Raw("select table_name from information_schema.tables where table_schema = ?", m.schema).Scan(&tables).Error
	if err != nil {
		return nil, err
	}
	return tables, nil
}

// FindColumns 查询列
func (m *PostgresqlModel) FindColumns(table string) ([]*DbColumn, error) {
	sql := "SELECT t.num,t.field AS column_name,t.type AS data_type,t.not_null AS is_nullable,t.comment,c.column_default FROM(SELECT a.attnum AS num,c.relname,a.attname AS field,t.typname AS TYPE,a.atttypmod AS lengthvar,a.attnotnull AS not_null,b.description AS COMMENT FROM pg_class c,pg_attribute a LEFT OUTER JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum=b.objsubid,pg_type t WHERE c.relname= ? AND a.attnum>0 AND a.attrelid=c.oid AND a.atttypid=t.oid GROUP BY a.attnum,c.relname,a.attname,t.typname,a.atttypmod,a.attnotnull,b.description ORDER BY a.attnum)AS t LEFT JOIN information_schema.columns AS c ON t.relname=c.table_name AND t.field=c.column_name AND c.table_schema= ?"
	var reply []*DbColumn
	err := m.gorm.Raw(sql, table, m.schema).Scan(&reply).Error
	if err != nil {
		return nil, err
	}
	return reply, nil
}

// FindIndex 查询索引
func (m *PostgresqlModel) FindIndex(table string) ([]*DbIndex, error) {
	querySql := `SELECT ix.relname as index_name, upper(am.amname) AS index_algorithm, indisunique as is_unique, pg_get_indexdef(indexrelid) as index_definition, replace(regexp_replace(regexp_replace(regexp_replace(pg_get_indexdef(indexrelid), ' WHERE .+|INCLUDE .+', ''), ' WITH .+', ''), '.*\((.*)\)', '\1'), ' ', '') AS column_name, CASE WHEN position(' WHERE 'in pg_get_indexdef(indexrelid))>0 THEN regexp_replace(pg_get_indexdef(indexrelid),'.+WHERE ','') WHEN position(' WITH 'in pg_get_indexdef(indexrelid))>0 THEN regexp_replace(pg_get_indexdef(indexrelid),'.+WITH ','') ELSE''END AS condition,pg_catalog.obj_description(i.indexrelid,'pg_class')as comment FROM pg_index i JOIN pg_class t ON t.oid = i.indrelid JOIN pg_class ix ON ix.oid = i.indexrelid JOIN pg_namespace n ON t.relnamespace = n.oid JOIN pg_am as am ON ix.relam = am.oid WHERE t.relname = ? AND n.nspname = ? `
	reply := make([]*DbIndex, 0)
	err := m.gorm.Raw(querySql, table, m.schema).Scan(&reply).Error
	if err != nil {
		return nil, err
	}
	return reply, nil
}
