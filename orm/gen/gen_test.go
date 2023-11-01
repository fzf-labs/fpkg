package gen

import (
	"testing"

	"github.com/fzf-labs/fpkg/orm"
	"github.com/stretchr/testify/assert"
	"gorm.io/gen"
	"gorm.io/gen/field"
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
	NewGenerationDB(client, "./example/postgres/", WithGenerateModel(func(g *gen.Generator) map[string]any {
		adminLogDemo := g.GenerateModel("admin_log_demo")
		AdminRoleDemo := g.GenerateModel("admin_role_demo",
			gen.FieldRelate(field.Many2Many, "Admins", g.GenerateModel("admin_demo"),
				&field.RelateConfig{
					RelateSlicePointer: true,
					JSONTag:            JSONTagNameStrategy("Admins"),
					GORMTag:            field.GormTag{"joinForeignKey": []string{"role_id"}, "joinReferences": []string{"admin_id"}, "many2many": []string{"admin_to_role_demo"}},
				},
			),
		)
		adminDemo := g.GenerateModel("admin_demo",
			gen.FieldRelate(field.HasMany, "AdminLogDemos", adminLogDemo,
				&field.RelateConfig{
					RelateSlicePointer: true,
					JSONTag:            JSONTagNameStrategy("AdminLogDemos"),
					GORMTag:            field.GormTag{"foreignKey": []string{"admin_id"}},
				},
			),
			gen.FieldRelate(field.Many2Many, "AdminRoles", AdminRoleDemo,
				&field.RelateConfig{
					RelateSlicePointer: true,
					JSONTag:            JSONTagNameStrategy("AdminRoles"),
					GORMTag:            field.GormTag{"joinForeignKey": []string{"admin_id"}, "joinReferences": []string{"role_id"}, "many2many": []string{"admin_to_role_demo"}},
				},
			),
		)
		return map[string]any{
			"admin_demo":     adminDemo,
			"admin_log_demo": adminLogDemo,
		}
	}), WithDataMap(DataTypeMap()), WithDBOpts(ModelOptionRemoveDefault(), ModelOptionUnderline("ul_"))).Do()
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
	NewGenerationDB(client, "./example/postgres/", WithDataMap(DataTypeMap()), WithDBOpts(ModelOptionRemoveDefault(), ModelOptionUnderline("UL"))).Do()
	assert.Equal(t, nil, err)
}

func TestNewGenerationPB(t *testing.T) {
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
	NewGenerationPB(client, "./example/postgres/pb", "api.gorm_gen.v1", "api/gorm_gen/v1;v1", WithPBOpts(ModelOptionRemoveDefault(), ModelOptionUnderline("ul_"))).Do()
	assert.Equal(t, nil, err)
}
