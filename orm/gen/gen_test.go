package gen

import (
	"testing"

	"github.com/fzf-labs/fpkg/orm"
	"github.com/stretchr/testify/assert"
)

func TestGenerationPostgres(t *testing.T) {
	client, err := orm.NewGormPostgresClient(&orm.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=gorm_gen sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	})
	if err != nil {
		return
	}
	NewGenerationDB(client, "./example/postgres/", WithDataMap(DefaultPostgresDataMap), WithDBOpts(ModelOptionRemoveDefault(), ModelOptionUnderline("ul_"))).Do()
	assert.Equal(t, nil, err)
}

func TestGenerationMysql(t *testing.T) {
	client, err := orm.NewGormMysqlClient(&orm.GormMysqlClientConfig{
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
	NewGenerationDB(client, "./example/postgres/", WithDataMap(DefaultPostgresDataMap), WithDBOpts(ModelOptionRemoveDefault(), ModelOptionUnderline("UL"))).Do()
	assert.Equal(t, nil, err)
}

func TestNewGenerationPb(t *testing.T) {
	client, err := orm.NewGormPostgresClient(&orm.GormPostgresClientConfig{
		DataSourceName:  "host=0.0.0.0 port=5432 user=postgres password=123456 dbname=gorm_gen sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConn:     0,
		MaxOpenConn:     0,
		ConnMaxLifeTime: 0,
		ShowLog:         false,
		Tracing:         false,
	})
	if err != nil {
		return
	}
	NewGenerationPb(client, "./example/postgres/pb", "api.gorm_gen.v1", "api/gorm_gen/v1;v1", WithPbOpts(ModelOptionRemoveDefault(), ModelOptionUnderline("ul_"))).Do()
	assert.Equal(t, nil, err)
}
