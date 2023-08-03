package repo

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
	"gorm.io/gorm"
)

type Repo struct {
	gorm         *gorm.DB
	daoPath      string
	modelPath    string
	relativePath string
}

func NewGenerationRepo(gorm *gorm.DB, daoPath string, modelPath string, relativePath string) *Repo {
	return &Repo{
		gorm:         gorm,
		daoPath:      daoPath,
		modelPath:    modelPath,
		relativePath: relativePath,
	}
}

func (r *Repo) MkdirPath() error {
	if err := os.MkdirAll(r.relativePath, os.ModePerm); err != nil {
		return fmt.Errorf("create model pkg path(%s) fail: %s", r.relativePath, err)
	}
	return nil
}

func (r *Repo) GenerationTable(table string, columnNameToDataType map[string]string) error {
	var file string

	var createMethods string
	var updateMethods string
	var findMethods string
	var delMethods string

	var createFunc string
	var updateFunc string
	var findFunc string
	var delFunc string

	dbName := r.gorm.Migrator().CurrentDatabase()
	indexes, err := r.gorm.Migrator().GetIndexes(table)
	if err != nil {
		return err
	}
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i].Name() > indexes[j].Name()
	})
	pkgTpl, err := NewTemplate("Pkg").Parse(Pkg).Execute(map[string]any{
		"lowerDbName": dbName,
	})
	if err != nil {
		return err
	}
	importTpl, err := NewTemplate("Import").Parse(Import).Execute(map[string]any{
		"FillDaoPkgPath":   FillModelPkgPath(r.daoPath),
		"FillModelPkgPath": FillModelPkgPath(r.modelPath),
		"relativePath":     r.relativePath,
		"lowerDbName":      dbName,
	})
	if err != nil {
		return err
	}
	upperTableName := r.UpperName(table)
	lowerTableName := r.LowerName(table)
	interfaceCreateOneTpl, err := NewTemplate("InterfaceCreateOne").Parse(InterfaceCreateOne).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
	})
	if err != nil {
		return err
	}
	createOneTpl, err := NewTemplate("CreateOne").Parse(CreateOne).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
	})
	if err != nil {
		return err
	}
	interfaceUpdateOneTpl, err := NewTemplate("InterfaceUpdateOne").Parse(InterfaceUpdateOne).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
	})
	if err != nil {
		return err
	}
	updateOneTpl, err := NewTemplate("UpdateOne").Parse(UpdateOne).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
	})
	if err != nil {
		return err
	}
	createMethods += fmt.Sprintln(interfaceCreateOneTpl.String())
	updateMethods += fmt.Sprintln(interfaceUpdateOneTpl.String())
	createFunc += fmt.Sprintln(createOneTpl.String())
	updateFunc += fmt.Sprintln(updateOneTpl.String())
	var cacheKeys string
	var varSingleCache string
	var varSingleCacheDel string
	for _, index := range indexes {
		//唯一性索引
		unique, _ := index.Unique()
		if unique {
			var cacheField string
			cacheFieldsJoinSli := make([]string, 0)
			for _, column := range index.Columns() {
				cacheField += r.UpperName(column)
				cacheFieldsJoinSli = append(cacheFieldsJoinSli, fmt.Sprintf("v.%s", r.UpperName(column)))
			}
			varCacheTpl, err := NewTemplate("VarCache").Parse(VarCache).Execute(map[string]any{
				"upperTableName": upperTableName,
				"cacheField":     cacheField,
			})
			if err != nil {
				return err
			}
			varSingleCacheTpl, err := NewTemplate("VarSingleCache").Parse(VarSingleCache).Execute(map[string]any{
				"upperTableName": upperTableName,
				"cacheField":     cacheField,
			})
			if err != nil {
				return err
			}
			varSingleCacheDelTpl, err := NewTemplate("VarSingleCacheDel").Parse(VarSingleCacheDel).Execute(map[string]any{
				"upperTableName":  upperTableName,
				"cacheField":      cacheField,
				"cacheFieldsJoin": strings.Join(cacheFieldsJoinSli, ","),
			})
			if err != nil {
				return err
			}
			cacheKeys += varCacheTpl.String()
			varSingleCache += fmt.Sprintln(varSingleCacheTpl.String())
			varSingleCacheDel += fmt.Sprintln(varSingleCacheDelTpl.String())
		}

		if len(index.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			fieldsJoin := strings.Join(index.Columns(), ",")
			var whereFields string
			for _, v := range index.Columns() {
				upperFields += r.UpperName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerName(v), columnNameToDataType[v])
				whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperName(v), r.LowerName(v))
			}
			if unique {
				interfaceFindOneCacheByFields, err := NewTemplate("InterfaceFindOneCacheByFields").Parse(InterfaceFindOneCacheByFields).Execute(map[string]any{
					"lowerDbName":       dbName,
					"upperTableName":    upperTableName,
					"lowerTableName":    lowerTableName,
					"upperFields":       upperFields,
					"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				})
				if err != nil {
					return err
				}
				findMethods += fmt.Sprintln(interfaceFindOneCacheByFields.String())
				findOneCacheByFields, err := NewTemplate("FindOneCacheByFields").Parse(FindOneCacheByFields).Execute(map[string]any{
					"lowerDbName":       dbName,
					"upperTableName":    upperTableName,
					"lowerTableName":    lowerTableName,
					"upperFields":       upperFields,
					"fieldsJoin":        fieldsJoin,
					"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
					"whereFields":       strings.Trim(whereFields, ","),
				})
				if err != nil {
					return err
				}
				findFunc += fmt.Sprintln(findOneCacheByFields.String())
				interfaceDeleteOneByFields, err := NewTemplate("InterfaceDeleteOneCacheByFields").Parse(InterfaceDeleteOneCacheByFields).Execute(map[string]any{
					"lowerDbName":       dbName,
					"upperTableName":    upperTableName,
					"lowerTableName":    lowerTableName,
					"upperField":        r.UpperName(index.Columns()[0]),
					"lowerField":        r.LowerName(index.Columns()[0]),
					"upperFields":       upperFields,
					"dataType":          columnNameToDataType[index.Columns()[0]],
					"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
					"whereFields":       strings.Trim(whereFields, ","),
				})
				if err != nil {
					return err
				}
				deleteOneByFields, err := NewTemplate("DeleteOneCacheByFields").Parse(DeleteOneCacheByFields).Execute(map[string]any{
					"lowerDbName":       dbName,
					"upperTableName":    upperTableName,
					"lowerTableName":    lowerTableName,
					"upperField":        r.UpperName(index.Columns()[0]),
					"lowerField":        r.LowerName(index.Columns()[0]),
					"upperFields":       upperFields,
					"dataType":          columnNameToDataType[index.Columns()[0]],
					"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
					"whereFields":       strings.Trim(whereFields, ","),
					"fieldsJoin":        fieldsJoin,
				})
				if err != nil {
					return err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneByFields.String())
				delFunc += fmt.Sprintln(deleteOneByFields.String())
			} else {
				interfaceFindMultiByFields, err := NewTemplate("InterfaceFindMultiByFields").Parse(InterfaceFindMultiByFields).Execute(map[string]any{
					"lowerDbName":       dbName,
					"upperTableName":    upperTableName,
					"lowerTableName":    lowerTableName,
					"upperFields":       upperFields,
					"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				})
				if err != nil {
					return err
				}
				findMethods += fmt.Sprintln(interfaceFindMultiByFields.String())
				findMultiByFields, err := NewTemplate("FindMultiByFields").Parse(FindMultiByFields).Execute(map[string]any{
					"lowerDbName":       dbName,
					"upperTableName":    upperTableName,
					"lowerTableName":    lowerTableName,
					"upperFields":       upperFields,
					"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
					"whereFields":       strings.Trim(whereFields, ","),
				})
				if err != nil {
					return err
				}
				findFunc += fmt.Sprintln(findMultiByFields.String())
			}

		} else {
			if unique {
				interfaceFindOneCacheByField, err := NewTemplate("InterfaceFindOneCacheByField").Parse(InterfaceFindOneCacheByField).Execute(map[string]any{
					"lowerDbName": dbName,

					"upperTableName": upperTableName,
					"lowerTableName": lowerTableName,
					"upperField":     r.UpperName(index.Columns()[0]),
					"lowerField":     r.LowerName(index.Columns()[0]),
					"dataType":       columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				findOneCacheByField, err := NewTemplate("FindOneCacheByField").Parse(FindOneCacheByField).Execute(map[string]any{
					"lowerDbName": dbName,

					"upperTableName": upperTableName,
					"lowerTableName": lowerTableName,
					"upperField":     r.UpperName(index.Columns()[0]),
					"lowerField":     r.LowerName(index.Columns()[0]),
					"dataType":       columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				findMethods += fmt.Sprintln(interfaceFindOneCacheByField.String())
				findFunc += fmt.Sprintln(findOneCacheByField.String())

				interfaceFindMultiCacheByFieldPlural, err := NewTemplate("InterfaceFindMultiCacheByFieldPlural").Parse(InterfaceFindMultiCacheByFieldPlural).Execute(map[string]any{
					"lowerDbName": dbName,

					"upperTableName":   upperTableName,
					"lowerTableName":   lowerTableName,
					"upperField":       r.UpperName(index.Columns()[0]),
					"lowerField":       r.LowerName(index.Columns()[0]),
					"upperFieldPlural": inflection.Plural(r.UpperName(index.Columns()[0])),
					"lowerFieldPlural": inflection.Plural(r.LowerName(index.Columns()[0])),
					"dataType":         columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				findMultiCacheByFieldPlural, err := NewTemplate("FindMultiCacheByFieldPlural").Parse(FindMultiCacheByFieldPlural).Execute(map[string]any{
					"lowerDbName": dbName,

					"upperTableName":   upperTableName,
					"lowerTableName":   lowerTableName,
					"upperField":       r.UpperName(index.Columns()[0]),
					"lowerField":       r.LowerName(index.Columns()[0]),
					"upperFieldPlural": inflection.Plural(r.UpperName(index.Columns()[0])),
					"lowerFieldPlural": inflection.Plural(r.LowerName(index.Columns()[0])),
					"dataType":         columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				findMethods += fmt.Sprintln(interfaceFindMultiCacheByFieldPlural.String())
				findFunc += fmt.Sprintln(findMultiCacheByFieldPlural.String())

				interfaceDeleteOneByField, err := NewTemplate("InterfaceDeleteOneCacheByField").Parse(InterfaceDeleteOneCacheByField).Execute(map[string]any{
					"lowerDbName":    dbName,
					"upperTableName": upperTableName,
					"lowerTableName": lowerTableName,
					"upperField":     r.UpperName(index.Columns()[0]),
					"lowerField":     r.LowerName(index.Columns()[0]),
					"dataType":       columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				deleteOneByField, err := NewTemplate("DeleteOneCacheByField").Parse(DeleteOneCacheByField).Execute(map[string]any{
					"lowerDbName":    dbName,
					"upperTableName": upperTableName,
					"lowerTableName": lowerTableName,
					"upperField":     r.UpperName(index.Columns()[0]),
					"lowerField":     r.LowerName(index.Columns()[0]),
					"dataType":       columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneByField.String())
				delFunc += fmt.Sprintln(deleteOneByField.String())
				interfaceDeleteMultiByFieldPlural, err := NewTemplate("InterfaceDeleteMultiCacheByFieldPlural").Parse(InterfaceDeleteMultiCacheByFieldPlural).Execute(map[string]any{
					"lowerDbName":      dbName,
					"upperTableName":   upperTableName,
					"lowerTableName":   lowerTableName,
					"upperField":       r.UpperName(index.Columns()[0]),
					"lowerField":       r.LowerName(index.Columns()[0]),
					"upperFieldPlural": inflection.Plural(r.UpperName(index.Columns()[0])),
					"lowerFieldPlural": inflection.Plural(r.LowerName(index.Columns()[0])),
					"dataType":         columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				deleteMultiByFieldPlural, err := NewTemplate("DeleteMultiCacheByFieldPlural").Parse(DeleteMultiCacheByFieldPlural).Execute(map[string]any{
					"lowerDbName":      dbName,
					"upperTableName":   upperTableName,
					"lowerTableName":   lowerTableName,
					"upperField":       r.UpperName(index.Columns()[0]),
					"lowerField":       r.LowerName(index.Columns()[0]),
					"upperFieldPlural": inflection.Plural(r.UpperName(index.Columns()[0])),
					"lowerFieldPlural": inflection.Plural(r.LowerName(index.Columns()[0])),
					"dataType":         columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiByFieldPlural.String())
				delFunc += fmt.Sprintln(deleteMultiByFieldPlural.String())
			} else {
				interfaceFindMultiByFieldPlural, err := NewTemplate("InterfaceFindMultiByFieldPlural").Parse(InterfaceFindMultiByFieldPlural).Execute(map[string]any{
					"lowerDbName":      dbName,
					"upperTableName":   upperTableName,
					"lowerTableName":   lowerTableName,
					"upperFieldPlural": inflection.Plural(r.UpperName(index.Columns()[0])),
					"lowerFieldPlural": inflection.Plural(r.LowerName(index.Columns()[0])),
					"dataType":         columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				findMultiByFieldPlural, err := NewTemplate("FindMultiByFieldPlural").Parse(FindMultiByFieldPlural).Execute(map[string]any{
					"lowerDbName":      dbName,
					"upperField":       r.UpperName(index.Columns()[0]),
					"upperTableName":   upperTableName,
					"lowerTableName":   lowerTableName,
					"upperFieldPlural": inflection.Plural(r.UpperName(index.Columns()[0])),
					"lowerFieldPlural": inflection.Plural(r.LowerName(index.Columns()[0])),
					"dataType":         columnNameToDataType[index.Columns()[0]],
				})
				if err != nil {
					return err
				}
				findMethods += fmt.Sprintln(interfaceFindMultiByFieldPlural.String())
				findFunc += fmt.Sprintln(findMultiByFieldPlural.String())
			}
		}
	}
	varTpl, err := NewTemplate("Var").Parse(Var).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
		"cacheKeys":      cacheKeys,
	})
	if err != nil {
		return err
	}
	newTpl, err := NewTemplate("New").Parse(New).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
	})
	if err != nil {
		return err
	}
	interfaceDeleteUniqueIndexCacheTpl, err := NewTemplate("InterfaceDeleteUniqueIndexCache").Parse(InterfaceDeleteUniqueIndexCache).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
		"singleCache":    varSingleCache,
		"singleCacheDel": varSingleCacheDel,
	})
	if err != nil {
		return err
	}
	deleteUniqueIndexCacheTpl, err := NewTemplate("DeleteUniqueIndexCache").Parse(DeleteUniqueIndexCache).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
		"singleCache":    varSingleCache,
		"singleCacheDel": varSingleCacheDel,
	})
	if err != nil {
		return err
	}
	delMethods += fmt.Sprintln(interfaceDeleteUniqueIndexCacheTpl.String())
	delFunc += fmt.Sprintln(deleteUniqueIndexCacheTpl.String())

	var methods string
	methods += createMethods
	methods += updateMethods
	methods += delMethods
	methods += findMethods
	typesTpl, err := NewTemplate("Types").Parse(Types).Execute(map[string]any{
		"lowerDbName":    dbName,
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
		"methods":        methods,
	})
	if err != nil {
		return err
	}
	file += fmt.Sprintln(pkgTpl.String())
	file += fmt.Sprintln(importTpl.String())
	file += fmt.Sprintln(varTpl.String())
	file += fmt.Sprintln(typesTpl.String())
	file += fmt.Sprintln(newTpl.String())
	file += fmt.Sprintln(createFunc)
	file += fmt.Sprintln(updateFunc)
	file += fmt.Sprintln(delFunc)
	file += fmt.Sprintln(findFunc)
	outputFile := r.relativePath + "/" + table + ".repo.go"
	err = r.output(outputFile, []byte(file))
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) output(fileName string, content []byte) error {
	result, err := imports.Process(fileName, content, nil)
	if err != nil {
		lines := strings.Split(string(content), "\n")
		errLine, _ := strconv.Atoi(strings.Split(err.Error(), ":")[1])
		startLine, endLine := errLine-5, errLine+5
		fmt.Println("Format fail:", errLine, err)
		if startLine < 0 {
			startLine = 0
		}
		if endLine > len(lines)-1 {
			endLine = len(lines) - 1
		}
		for i := startLine; i <= endLine; i++ {
			fmt.Println(i, lines[i])
		}
		return fmt.Errorf("cannot format file: %w", err)
	}
	return os.WriteFile(fileName, result, 0640)
}

func FillModelPkgPath(filePath string) string {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName,
		Dir:  filePath,
	})
	if err != nil {
		return ""
	}
	if len(pkgs) == 0 {
		return ""
	}
	return pkgs[0].PkgPath
}
func (r *Repo) UpperName(s string) string {
	return r.gorm.NamingStrategy.SchemaName(s)
}

func (r *Repo) LowerName(s string) string {
	s = r.UpperName(s)
	if len(s) == 0 {
		return s
	}
	commonInitialisms := []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	//如果第一个单词命中  则不处理
	for _, v := range commonInitialisms {
		if strings.HasPrefix(s, v) {
			return s
		}
	}
	rs := []rune(s)
	f := rs[0]

	if 'A' <= f && f <= 'Z' {
		return string(unicode.ToLower(f)) + string(rs[1:])
	}
	return s
}
