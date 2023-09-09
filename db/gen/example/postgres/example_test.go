package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/fzf-labs/fpkg/cache/rueidiscache"
	"github.com/fzf-labs/fpkg/db"
	rueidiscache2 "github.com/fzf-labs/fpkg/db/gen/cache/rueidisdbcache"
	"github.com/fzf-labs/fpkg/db/gen/example/postgres/gorm_gen_repo"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	gormPostgresClient, err := db.NewGormPostgresClient(&db.GormPostgresClientConfig{
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
	client, err := rueidiscache.NewRueidis(&rueidis.ClientOption{
		Username:    "",
		Password:    "123456",
		InitAddress: []string{"127.0.0.1:6379"},
		SelectDB:    0,
	})
	if err != nil {
		return
	}
	ctx := context.Background()
	rueidisCache := rueidiscache2.NewRueidisCache(client)
	repo := gorm_gen_repo.NewUserDemoRepo(gormPostgresClient, rueidisCache)
	result, err := repo.FindOneCacheByID(ctx, 1)
	if err != nil {
		return
	}
	fmt.Println(result)
	assert.Equal(t, nil, err)
}
