package gen

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
	"unicode"
)

func Generation(db *gorm.DB, dataMap map[string]func(detailType string) (dataType string), outPath string, modelPkgPath string) {
	// 初始化
	g := gen.NewGenerator(gen.Config{
		OutPath:      outPath,
		ModelPkgPath: modelPkgPath,
	})
	// 使用数据库
	g.UseDB(db)
	// 自定义字段类型映射
	g.WithDataTypeMap(dataMap)
	// json 小驼峰模型命名
	g.WithJSONTagNameStrategy(func(c string) string {
		return LowerCamelCase(c)
	})
	// 从数据库中生成所有表
	g.ApplyBasic(g.GenerateAllTable()...)
	// 在结构或表模型上应用diy接口
	//g.ApplyInterface(func(model.Method) {}, g.GenerateModel("user"))
	g.Execute()
}

// 默认mysql字段类型映射
var defaultMySqlDataMap = map[string]func(detailType string) (dataType string){
	"int":     func(detailType string) (dataType string) { return "int64" },
	"tinyint": func(detailType string) (dataType string) { return "int32" },
	"json":    func(string) string { return "datatypes.JSON" },
}

// 默认Postgres字段类型映射
var defaultPostgresDataMap = map[string]func(detailType string) (dataType string){
	"json": func(string) string { return "datatypes.JSON" },
}

// ConnectDB 数据库连接
func ConnectDB(dbType string, dsn string) *gorm.DB {
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

// UpperCamelCase 下划线单词转为大写驼峰单词
func UpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	causer := cases.Title(language.English)
	s = causer.String(s)
	return strings.Replace(s, " ", "", -1)
}

// LowerCamelCase 下划线单词转为小写驼峰单词
func LowerCamelCase(s string) string {
	s = UpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}
