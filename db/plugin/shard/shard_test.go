package shard

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"testing"
)

func TestNewMonthShardingPlugin(t *testing.T) {
	sqlDB, err := sql.Open("mysql", "root:123456@tcp(0.0.0.0:3306)/fgzs-single?charset=utf8mb4&loc=Asia%2FShanghai&parseTime=true")
	if err != nil {
		fmt.Printf("open mysql failed! err: %+v", err)
	}
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}
	gormConfig.Logger = logger.Default.LogMode(logger.Info)
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gormConfig)
	if err != nil {
		fmt.Printf("database connection failed!  err: %+v", err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	err = db.Use(NewMonthShardingPlugin("orders", "time"))
	if err != nil {
		fmt.Printf("gormopentracing new failed!  err: %+v", err)
	}
	// this record will insert to orders_03
	err = db.Exec("INSERT INTO orders(time) VALUES(?)", int64(1669097233)).Error
	if err != nil {
		fmt.Println(err)
	}

}
