package gen

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/fzf-labs/fpkg/orm/gen/repo"
	"github.com/iancoleman/strcase"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const (
	SQLNullTime = "sql.NullTime"
	TimeTime    = "time.Time"
)

type Generation struct {
	db         *gorm.DB                                                      // 数据库
	outPutPath string                                                        // 文件生成链接
	genRepo    bool                                                          // 是否生成repo
	dataMap    map[string]func(columnType gorm.ColumnType) (dataType string) // 自定义字段类型映射
	tables     []string                                                      // 指定表集合
	opts       []gen.ModelOpt                                                // 特殊处理逻辑函数
}

func NewGeneration(db *gorm.DB, outPutPath string, opts ...Option) *Generation {
	g := &Generation{
		db:         db,
		outPutPath: outPutPath,
		genRepo:    true,
		dataMap:    nil,
		tables:     nil,
		opts:       nil,
	}
	if len(opts) > 0 {
		for _, v := range opts {
			v(g)
		}
	}
	return g
}

type Option func(gen *Generation)

// WithOutRepo 选项函数-不生成repo
func WithOutRepo() Option {
	return func(r *Generation) {
		r.genRepo = false
	}
}

// WithTables 选项函数-自定义表
func WithTables(tables []string) Option {
	return func(r *Generation) {
		r.tables = tables
	}
}

// WithDataMap 选项函数-自定义关系映射
func WithDataMap(dataMap map[string]func(columnType gorm.ColumnType) (dataType string)) Option {
	return func(r *Generation) {
		r.dataMap = dataMap
	}
}

// WithOpts 选项函数-自定义特殊设置
func WithOpts(opts ...gen.ModelOpt) Option {
	return func(r *Generation) {
		r.opts = opts
	}
}

// Do 生成
func (g *Generation) Do() {
	// 路径处理
	dbName := g.db.Migrator().CurrentDatabase()
	outPutPath := strings.Trim(g.outPutPath, "/")
	daoPath := fmt.Sprintf("%s/%s_dao", outPutPath, dbName)
	modelPath := fmt.Sprintf("%s/%s_model", outPutPath, dbName)
	repoPath := fmt.Sprintf("%s/%s_repo", outPutPath, dbName)
	// 初始化
	generator := gen.NewGenerator(gen.Config{
		OutPath:      daoPath,
		ModelPkgPath: modelPath,
	})
	// 使用数据库
	generator.UseDB(g.db)
	// 自定义字段类型映射
	generator.WithDataTypeMap(g.dataMap)
	// json 小驼峰模型命名
	generator.WithJSONTagNameStrategy(LowerCamelCase)
	// 特殊处理逻辑
	if len(g.opts) > 0 {
		generator.WithOpts(g.opts...)
	}
	// 获取所有表
	tables, err := g.db.Migrator().GetTables()
	if err != nil {
		return
	}
	// 指定表
	if len(g.tables) > 0 {
		tables = g.tables
	}
	models := make([]any, len(tables))
	for i, tableName := range tables {
		models[i] = generator.GenerateModel(tableName)
	}
	generator.ApplyBasic(models...)
	// 生成model,dao
	generator.Execute()
	// 判断是否生成repo
	if !g.genRepo {
		return
	}
	// 生成repo
	generationRepo := repo.NewGenerationRepo(g.db, daoPath, modelPath, repoPath)
	// 生成repo的文件夹目录文件
	err = generationRepo.MkdirPath()
	if err != nil {
		log.Println("repo MkdirPath err:", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(tables))
	for _, v := range tables {
		t := v
		go func(table string) {
			defer wg.Done()
			// 表字段对应的类型
			columnNameToDataType := make(map[string]string)
			// 表字段对应的名称
			columnNameToName := make(map[string]string)
			// 表字段对应的dao字段类型
			columnNameToFieldType := make(map[string]string)
			queryStructMeta := generator.GenerateModel(table)
			for _, vv := range queryStructMeta.Fields {
				columnNameToDataType[vv.ColumnName] = vv.Type
				columnNameToName[vv.ColumnName] = vv.Name
				columnNameToFieldType[vv.ColumnName] = vv.GenType()
			}
			// 数据表repo代码生成
			err = generationRepo.GenerationTable(table, columnNameToDataType, columnNameToName, columnNameToFieldType)
			if err != nil {
				log.Println("repo GenerationTable err:", err)
				return
			}
		}(t)
	}
	wg.Wait()
}

// ModelOptionUnderline 前缀是下划线重命名
func ModelOptionUnderline(rename string) gen.ModelOpt {
	return gen.FieldModify(func(f gen.Field) gen.Field {
		if strings.HasPrefix(f.Name, "_") {
			f.Name = strings.Replace(f.Name, "_", rename, 1)
			f.Tag.Set(field.TagKeyJson, f.ColumnName)
		}
		return f
	})
}

// ModelOptionPgDefaultString Postgres默认字符串处理
func ModelOptionPgDefaultString() gen.ModelOpt {
	return gen.FieldGORMTagReg(".*?", func(tag field.GormTag) field.GormTag {
		regex := regexp.MustCompile(`default:'(.*?)'::character varying`)
		matches := regex.FindStringSubmatch(tag.Build())
		if len(matches) > 0 {
			tag.Set("default", matches[1])
		}
		return tag
	})
}

func ModelOptionRemoveDefault() gen.ModelOpt {
	return gen.FieldGORMTagReg(".*?", func(tag field.GormTag) field.GormTag {
		tag.Remove("default")
		return tag
	})
}

// DefaultPostgresDataMap 默认Postgres字段类型映射
var DefaultPostgresDataMap = map[string]func(columnType gorm.ColumnType) (dataType string){
	"json":  func(columnType gorm.ColumnType) string { return "datatypes.JSON" },
	"jsonb": func(columnType gorm.ColumnType) string { return "datatypes.JSON" },
	"timestamptz": func(columnType gorm.ColumnType) string {
		if repo.StrSliFind([]string{"deleted_at", "deletedAt", "deleted_time", "deleted_time"}, columnType.Name()) {
			return "gorm.DeletedAt"
		}
		nullable, _ := columnType.Nullable()
		if nullable {
			return SQLNullTime
		}
		return TimeTime
	},
}

// ConnectDB 数据库连接
func ConnectDB(dbType, dsn string) *gorm.DB {
	var db *gorm.DB
	var err error
	switch dbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			panic(fmt.Errorf("connect mysql db fail: %s", err))
		}
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn))
		if err != nil {
			panic(fmt.Errorf("connect postgres db fail: %s", err))
		}
	default:
		panic(fmt.Errorf(" db type err"))
	}
	return db
}

// LowerCamelCase 下划线单词转为小写驼峰单词
func LowerCamelCase(s string) string {
	return strcase.ToLowerCamel(s)
}
