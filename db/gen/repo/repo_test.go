package repo

import (
	"testing"

	"github.com/fzf-labs/fpkg/db"
)

func TestGenerationTable(t *testing.T) {
	client, _ := db.NewGormPostgresClient(&db.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	})
	newRepo := NewRepo(client, "github.com/fzf-labs/fpkg", "db/gen/example/postgres")
	err := newRepo.GenerationTable("system_users")
	if err != nil {
		return
	}
}
