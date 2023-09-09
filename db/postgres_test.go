package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, nil, err)
}

func TestDsnParse(t *testing.T) {
	parse := PostgresDsnParse("host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai")
	fmt.Println(parse)
	assert.Equal(t, &PostgresDsn{
		Host:     "0.0.0.0",
		Port:     5432,
		User:     "postgres",
		Password: "123456",
		Dbname:   "user",
	}, parse)
}

func TestDumpSQL(t *testing.T) {
	config := GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=gorm_gen sslmode=disable TimeZone=Asia/Shanghai",
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
	DumpPostgres(db, config.DataSourceName, "./gen/example/postgres/sql")
	assert.Equal(t, nil, err)
}
