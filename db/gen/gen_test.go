package gen

import (
	"testing"

	"github.com/fzf-labs/fpkg/db"
)

func TestGenerationPostgres(t *testing.T) {
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
	Generation(client, DefaultMySqlDataMap, "github.com/fzf-labs/fpkg", "./example/postgres/user_dao", "./example/postgres/user_model", "./example/postgres/user_repo")
}

func TestGenerationMysql(t *testing.T) {
	client, err := db.NewGormMysqlClient(&db.GormMysqlClientConfig{
		DataSourceName:  "",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	})
	if err != nil {
		return
	}
	Generation(client, DefaultPostgresDataMap, "", "./example/mysql/dao", "./example/mysql/model", "./example/postgres/user_repo")
}
