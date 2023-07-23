package db

import (
	"fmt"
	"testing"
)

func TestNewGormPostgresClient(t *testing.T) {
	config := GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	}
	_, err := NewGormPostgresClient(&config)
	fmt.Println(err)
	if err != nil {
		return
	}
}

func TestDsnParse(t *testing.T) {
	parse := DsnParse("host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai")
	fmt.Println(parse)
}

func TestDumpSql(t *testing.T) {
	config := GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	}
	db, err := NewGormPostgresClient(&config)
	fmt.Println(err)
	if err != nil {
		return
	}
	DumpSql(db, config.DataSourceName, "../sql")
}
