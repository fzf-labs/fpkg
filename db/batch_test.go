package db

import (
	"fmt"
	"testing"
)

type Demo struct {
	ID    int64   `gorm:"column:id;primary_key" json:"id"`
	Name  string  `gorm:"column:name" json:"name"`
	Width float64 `gorm:"column:width" json:"width"`
}

func TestBatchUpdateToSqlArray(t *testing.T) {
	var demo []*Demo
	demo1 := &Demo{1, "nihao", 12.1}
	demo2 := &Demo{2, "renzhen", 13.0}
	demo3 := &Demo{3, "duidai", 12.0}
	demo4 := &Demo{4, "xiexie", 13.0}
	demo5 := &Demo{5, "OOP", 12.0}
	demo = append(demo, demo1, demo2, demo3, demo4, demo5)
	sqls, err := BatchUpdateToSqlArray("demo_tab", demo)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(sqls)
}
