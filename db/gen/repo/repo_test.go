package repo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fzf-labs/fpkg/db"
)

func TestGenerationRepo(t *testing.T) {
	client, err := db.NewGormPostgresClient(&db.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	})
	if err != nil {
		return
	}
	fmt.Println(client.Migrator().CurrentDatabase())
	fmt.Println(client.Migrator().GetTables())
	indexes, err := client.Migrator().GetIndexes("user")
	if err != nil {
		return
	}
	marshal, err := json.Marshal(indexes)
	if err != nil {
		return
	}
	fmt.Println(string(marshal))
}
