package dbfunc

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	Postgres = "postgres"
	Mysql    = "mysql"
)

// SortIndexColumns 排序索引字段
func SortIndexColumns(db *gorm.DB, table string) (resp map[string][]string, err error) {
	switch db.Dialector.Name() {
	case Postgres:
		resp, err = pgSortIndexColumns(db, table)
		if err != nil {
			return nil, err
		}
	default:
	}
	return resp, nil
}

// pgSortIndexColumns  postgres索引字段排序
func pgSortIndexColumns(db *gorm.DB, table string) (map[string][]string, error) {
	resp := make(map[string][]string)
	type Tmp struct {
		TableName  string `gorm:"column:table_name" json:"table_name"`
		IndexName  string `gorm:"column:index_name" json:"index_name"`
		ColumnName string `gorm:"column:column_name" json:"column_name"`
	}
	result := make([]Tmp, 0)
	sql := fmt.Sprintf(`SELECT t.relname AS table_name,i.relname AS index_name,a.attname AS column_name,ix.indisunique AS non_unique,ix.indisprimary AS PRIMARY FROM pg_class t JOIN pg_index ix ON t.oid=ix.indrelid JOIN pg_class i ON i.oid=ix.indexrelid JOIN pg_attribute a ON a.attrelid=t.oid AND a.attnum=ANY(ix.indkey)WHERE t.relkind='r' AND t.relname='%s' ORDER BY ix.indrelid,(array_position(ix.indkey,a.attnum))`, table)
	err := db.Raw(sql).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		if _, ok := resp[v.IndexName]; !ok {
			resp[v.IndexName] = make([]string, 0)
		}
		resp[v.IndexName] = append(resp[v.IndexName], v.ColumnName)
	}
	return resp, nil
}

// GetPartitionTableName 获取分区表
func GetPartitionTableName(db *gorm.DB) (resp []string, err error) {
	switch db.Dialector.Name() {
	case Postgres:
		resp, err = getPGPartitionTableName(db)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(" db type err")
	}
	return resp, nil
}

// getPGPartitionTableName 获取分区表
func getPGPartitionTableName(db *gorm.DB) ([]string, error) {
	result := make([]string, 0)
	sql := `SELECT c.relname AS partitioned_table FROM pg_catalog.pg_class c JOIN pg_catalog.pg_inherits i ON c.oid=i.inhparent GROUP BY c.relname`
	err := db.Raw(sql).Pluck("partitioned_table", &result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetPartitionChildTableForTable 获取PG分区表的子表
func GetPartitionChildTableForTable(db *gorm.DB, tableName string) (resp []string, err error) {
	switch db.Dialector.Name() {
	case Postgres:
		resp, err = getPGPartitionChildTableForTable(db, tableName)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(" db type err")
	}
	return resp, nil
}

// getPGPartitionChildTableForTable 获取PG分区表的子表
func getPGPartitionChildTableForTable(db *gorm.DB, tableName string) ([]string, error) {
	result := make([]string, 0)
	sql := fmt.Sprintf(`SELECT c.relname AS child_table FROM pg_catalog.pg_class c JOIN pg_catalog.pg_inherits i ON c.oid=i.inhrelid WHERE i.inhparent=(SELECT oid FROM pg_class WHERE relname='%s')`, tableName)
	err := db.Raw(sql).Pluck("child_table", &result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetPartitionChildTable 获取所有分区表的子表
func GetPartitionChildTable(db *gorm.DB) (resp []string, err error) {
	switch db.Dialector.Name() {
	case Postgres:
		resp, err = getPartitionChildTable(db)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf(" db type err")
	}
	return resp, nil
}

// getPartitionChildTable 获取PG获取所有分区表的子表
func getPartitionChildTable(db *gorm.DB) ([]string, error) {
	result := make([]string, 0)
	sql := `SELECT c.relname AS child_table FROM pg_catalog.pg_class c JOIN pg_catalog.pg_inherits i ON c.oid=i.inhrelid`
	err := db.Raw(sql).Pluck("child_table", &result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
