package plugin

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestNewMonthShardingPlugin(t *testing.T) {
	sqlDB, err := sql.Open("pgx", "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=fkratos_sys sslmode=disable TimeZone=Asia/Shanghai")
	if err != nil {
		fmt.Printf("open mysql failed! err: %+v", err)
	}
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}
	gormConfig.Logger = logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gormConfig)
	if err != nil {
		fmt.Printf("database connection failed!  err: %+v", err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	err = db.Use(NewMonthShardingPlugin("sys_admin", "created_at"))
	if err != nil {
		fmt.Printf("gormopentracing new failed!  err: %+v", err)
	}
	// this record will insert to orders_03
	err = db.Exec("SELECT * FROM sys_admin WHERE created_at in ('2023-01-13 20:58:35')  ").Error
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, nil, err)
}

func TestNewShardingPlugin(t *testing.T) {
	sqlDB, err := sql.Open("pgx", "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=fkratos_sys sslmode=disable TimeZone=Asia/Shanghai")
	if err != nil {
		fmt.Printf("open mysql failed! err: %+v", err)
	}
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}
	gormConfig.Logger = logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gormConfig)
	if err != nil {
		fmt.Printf("database connection failed!  err: %+v", err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	err = db.Use(NewShardingPlugin("sys_admin", "created_at", 64))
	if err != nil {
		fmt.Printf("gormopentracing new failed!  err: %+v", err)
	}
	// this record will insert to orders_03
	err = db.Exec("SELECT * FROM sys_admin WHERE created_at  ='2023-01-13 20:58:01'  ").Error
	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, nil, err)
}
