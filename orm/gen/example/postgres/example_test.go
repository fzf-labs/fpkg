package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/fzf-labs/fpkg/cache/rueidiscache"
	"github.com/fzf-labs/fpkg/orm"
	"github.com/fzf-labs/fpkg/orm/gen/cache/rueidisdbcache"
	"github.com/fzf-labs/fpkg/orm/gen/example/postgres/gorm_gen_repo"
	"github.com/fzf-labs/fpkg/util/jsonutil"
	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
)

func Test_main(t *testing.T) {
	gormPostgresClient, err := orm.NewGormPostgresClient(&orm.GormPostgresClientConfig{
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
	rueidisCache := rueidisdbcache.NewRueidisDBCache(client)
	repo := gorm_gen_repo.NewUserDemoRepo(gormPostgresClient, rueidisCache)
	result, err := repo.FindOneCacheByID(ctx, 1)
	if err != nil {
		return
	}
	fmt.Println(result)
	assert.Equal(t, nil, err)
}

func Test_FindMultiByPaginator(t *testing.T) {
	gormPostgresClient, err := orm.NewGormPostgresClient(&orm.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=user sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         true,
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
	rueidisCache := rueidisdbcache.NewRueidisDBCache(client)
	repo := gorm_gen_repo.NewUserDemoRepo(gormPostgresClient, rueidisCache)
	result, total, err := repo.FindMultiByPaginator(ctx, &orm.PaginatorReq{
		Page:     1,
		PageSize: 1,
		Order:    "created_at,DESC",
		Search: []*orm.SearchColumn{
			{
				Field: "nickname",
				Value: "123",
				Exp:   "=",
				Logic: "",
			},
			//{
			//	Field: "status",
			//	Value: "0,-1",
			//	Exp:   "IN",
			//	Logic: "",
			//},
		},
	})
	jsonutil.Dump(result)
	fmt.Println(total)
	assert.Equal(t, nil, err)
}
