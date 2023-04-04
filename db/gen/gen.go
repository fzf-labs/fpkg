package gen

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
	"unicode"
)

type Config struct {
	OutPath      string
	ModelPkgPath string
}

func Generation(db *gorm.DB, conf Config) {
	g := gen.NewGenerator(gen.Config{
		OutPath:      conf.OutPath,
		ModelPkgPath: conf.ModelPkgPath,
	})
	g.UseDB(db)
	g.WithDataTypeMap(dataMap)
	g.WithJSONTagNameStrategy(func(c string) string {
		return LowerCamelCase(c)
	})
	// generate all table from database
	g.ApplyBasic(g.GenerateAllTable()...)
	// apply diy interfaces on structs or table models
	//g.ApplyInterface(func(model.Method) {}, g.GenerateModel("user"))
	g.Execute()
}

var dataMap = map[string]func(detailType string) (dataType string){
	"int":     func(detailType string) (dataType string) { return "int64" },
	"tinyint": func(detailType string) (dataType string) { return "int32" },
	"json":    func(string) string { return "datatypes.JSON" },
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
