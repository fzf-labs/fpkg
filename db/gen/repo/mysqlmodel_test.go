package repo

import (
	"fmt"
	"testing"

	"github.com/fzf-labs/fpkg/db"
	"github.com/fzf-labs/fpkg/util/jsonutil"
)

func TestMysqlModel_FindAllTables(t *testing.T) {
	config := db.GormMysqlClientConfig{
		DataSourceName:  "root:123456@tcp(0.0.0.0:3306)/user?charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	}
	gorm, err := db.NewGormMysqlClient(&config)
	if err != nil {
		return
	}
	model := NewMysqlModel(gorm)
	tables, err := model.FindAllTables()
	if err != nil {
		return
	}
	fmt.Println(tables)
}

func TestMysqlModel_FindColumns(t *testing.T) {
	config := db.GormMysqlClientConfig{
		DataSourceName:  "root:123456@tcp(0.0.0.0:3306)/user?charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	}
	gorm, err := db.NewGormMysqlClient(&config)
	if err != nil {
		return
	}
	model := NewMysqlModel(gorm)
	tables, err := model.FindColumns("user")
	if err != nil {
		return
	}
	jsonutil.Dump(tables)
}

func TestMysqlModel_FindIndex(t *testing.T) {
	config := db.GormMysqlClientConfig{
		DataSourceName:  "root:123456@tcp(0.0.0.0:3306)/user?charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	}
	gorm, err := db.NewGormMysqlClient(&config)
	if err != nil {
		return
	}
	model := NewMysqlModel(gorm)
	tables, err := model.FindIndex("user")
	if err != nil {
		return
	}
	jsonutil.Dump(tables)
}
