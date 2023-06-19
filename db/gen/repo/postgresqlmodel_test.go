package repo

import (
	"fmt"
	"testing"

	"github.com/fzf-labs/fpkg/db"
	"github.com/fzf-labs/fpkg/util/jsonutil"
)

func TestPostgresqlModel_GetAllTables(t *testing.T) {
	config := db.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	}
	gorm, err := db.NewGormPostgresClient(&config)
	if err != nil {
		return
	}
	newPostgresqlModel := NewPostgresqlModel(gorm)
	tables, err := newPostgresqlModel.FindAllTables()
	if err != nil {
		return
	}
	fmt.Println(tables)
}

func TestPostgresqlModel_FindIndex(t *testing.T) {
	config := db.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         true,
		Tracing:         false,
	}
	gorm, err := db.NewGormPostgresClient(&config)
	if err != nil {
		return
	}
	newPostgresqlModel := NewPostgresqlModel(gorm)
	tables, err := newPostgresqlModel.FindIndex("user")
	if err != nil {
		return
	}
	jsonutil.Dump(tables)
}

func TestPostgresqlModel_FindColumns(t *testing.T) {
	config := db.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         true,
		Tracing:         false,
	}
	gorm, err := db.NewGormPostgresClient(&config)
	if err != nil {
		return
	}
	newPostgresqlModel := NewPostgresqlModel(gorm)
	tables, err := newPostgresqlModel.FindColumns("user")
	if err != nil {
		return
	}
	jsonutil.Dump(tables)
}
