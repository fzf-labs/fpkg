package db

import (
	"fmt"
	"testing"
)

func TestNewGormPostgresClient(t *testing.T) {
	config := GormPostgresClientConfig{
		DataSourceName:  "",
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
