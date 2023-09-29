//nolint:all
package repo

import (
	"fmt"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
	"golang.org/x/tools/imports"
	"gorm.io/gorm"
)

var KeyWords = []string{
	"dao",
	"parameters",
	"cacheKey",
	"cacheKeys",
	"cacheValue",
	"keyToParam",
	"resp",
	"result",
	"marshal",
}

type Repo struct {
	gorm         *gorm.DB
	daoPath      string
	modelPath    string
	relativePath string
}

func NewGenerationRepo(db *gorm.DB, daoPath, modelPath, relativePath string) *Repo {
	return &Repo{
		gorm:         db,
		daoPath:      daoPath,
		modelPath:    modelPath,
		relativePath: relativePath,
	}
}

func (r *Repo) GenerationTable(table string, columnNameToDataType, columnNameToName, columnNameToFieldType map[string]string) error {
	var file string
	// 查询当前db的索引
	indexes, err := r.gorm.Migrator().GetIndexes(table)
	if err != nil {
		return err
	}
	indexes = r.ProcessIndex(indexes)
	lowerDBName := r.gorm.Migrator().CurrentDatabase()
	generationRepo := GenerationRepo{
		gorm:                  r.gorm,
		columnNameToDataType:  columnNameToDataType,
		columnNameToName:      columnNameToName,
		columnNameToFieldType: columnNameToFieldType,
		firstTableChar:        "",
		lowerDBName:           lowerDBName,
		lowerTableName:        "",
		upperTableName:        "",
		daoPkgPath:            FillModelPkgPath(r.daoPath),
		modelPkgPath:          FillModelPkgPath(r.modelPath),
		index:                 indexes,
	}
	generationRepo.lowerTableName = generationRepo.LowerName(table)
	generationRepo.upperTableName = generationRepo.UpperName(table)
	generationRepo.firstTableChar = generationRepo.lowerTableName[0:1]
	generatePkg, err := generationRepo.generatePkg()
	if err != nil {
		return err
	}
	generateImport, err := generationRepo.generateImport()
	if err != nil {
		return err
	}
	generateVar, err := generationRepo.generateVar()
	if err != nil {
		return err
	}
	generateTypes, err := generationRepo.generateTypes()
	if err != nil {
		return err
	}
	generateNew, err := generationRepo.generateNew()
	if err != nil {
		return err
	}
	generateCreateFunc, err := generationRepo.generateCreateFunc()
	if err != nil {
		return err
	}
	generateUpdateFunc, err := generationRepo.generateUpdateFunc()
	if err != nil {
		return err
	}
	generateDelFunc, err := generationRepo.generateDelFunc()
	if err != nil {
		return err
	}
	generateReadFunc, err := generationRepo.generateReadFunc()
	if err != nil {
		return err
	}
	file += fmt.Sprintln(generatePkg)
	file += fmt.Sprintln(generateImport)
	file += fmt.Sprintln(generateVar)
	file += fmt.Sprintln(generateTypes)
	file += fmt.Sprintln(generateNew)
	file += fmt.Sprintln(generateCreateFunc)
	file += fmt.Sprintln(generateUpdateFunc)
	file += fmt.Sprintln(generateDelFunc)
	file += fmt.Sprintln(generateReadFunc)
	outputFile := r.relativePath + "/" + table + ".repo.go"
	err = r.output(outputFile, []byte(file))
	if err != nil {
		return err
	}
	return nil
}

// MkdirPath 生成文件夹
func (r *Repo) MkdirPath() error {
	if err := os.MkdirAll(r.relativePath, os.ModePerm); err != nil {
		return fmt.Errorf("create model pkg path(%s) fail: %s", r.relativePath, err)
	}
	return nil
}

// ProcessIndex 索引处理  索引去重和排序
func (r *Repo) ProcessIndex(indexes []gorm.Index) []gorm.Index {
	repeat := make(map[string]struct{})
	result := make([]gorm.Index, 0)
	// 主键索引
	for _, v := range indexes {
		primaryKey, _ := v.PrimaryKey()
		if primaryKey {
			_, ok := repeat[SliToStr(v.Columns())]
			if !ok {
				repeat[SliToStr(v.Columns())] = struct{}{}
				result = append(result, v)
			}
		}
	}
	// 唯一索引
	for _, v := range indexes {
		primaryKey, _ := v.PrimaryKey()
		unique, _ := v.Unique()
		if !primaryKey && unique {
			_, ok := repeat[SliToStr(v.Columns())]
			if !ok {
				repeat[SliToStr(v.Columns())] = struct{}{}
				result = append(result, v)
			}
		}
	}
	// 普通索引
	for _, v := range indexes {
		primaryKey, _ := v.PrimaryKey()
		unique, _ := v.Unique()
		if !primaryKey && !unique {
			_, ok := repeat[SliToStr(v.Columns())]
			if !ok {
				repeat[SliToStr(v.Columns())] = struct{}{}
				result = append(result, v)
			}
		}
	}
	// 索引按名称排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name() > result[j].Name()
	})
	return result
}

// output 导出文件
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
	return os.WriteFile(fileName, result, 0600)
}

type GenerationRepo struct {
	gorm                  *gorm.DB
	columnNameToDataType  map[string]string // 字段名称对应的类型
	columnNameToName      map[string]string // 字段名称对应的Go名称
	columnNameToFieldType map[string]string // 字段名称对应的dao类型
	firstTableChar        string
	lowerDBName           string
	lowerTableName        string
	upperTableName        string
	daoPkgPath            string
	modelPkgPath          string
	index                 []gorm.Index
}

// generatePkg
func (r *GenerationRepo) generatePkg() (string, error) {
	tplParams := map[string]any{
		"lowerDBName": r.lowerDBName,
	}
	tpl, err := NewTemplate("Pkg").Parse(Pkg).Execute(tplParams)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

// generateImport
func (r *GenerationRepo) generateImport() (string, error) {
	tplParams := map[string]any{
		"daoPkgPath":   r.daoPkgPath,
		"modelPkgPath": r.daoPkgPath,
	}
	tpl, err := NewTemplate("Import").Parse(Import).Execute(tplParams)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

// generateVar
func (r *GenerationRepo) generateVar() (string, error) {
	var varStr string
	var cacheKeys string
	for _, v := range r.index {
		if r.CheckDaoFieldType(v.Columns()) {
			continue
		}
		unique, _ := v.Unique()
		if unique {
			var cacheField string
			for _, column := range v.Columns() {
				cacheField += r.UpperFieldName(column)
			}
			varCacheTpl, err := NewTemplate("VarCache").Parse(VarCache).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"cacheField":     cacheField,
			})
			if err != nil {
				return "", err
			}
			cacheKeys += varCacheTpl.String()
		}
	}
	varTpl, err := NewTemplate("Var").Parse(Var).Execute(map[string]any{
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	varStr += fmt.Sprintln(varTpl.String())
	if len(cacheKeys) > 0 {
		varCacheKeysTpl, err := NewTemplate("Var").Parse(VarCacheKeys).Execute(map[string]any{
			"cacheKeys": cacheKeys,
		})
		if err != nil {
			return "", err
		}
		varStr += fmt.Sprintln(varCacheKeysTpl.String())
	}
	return varStr, nil
}

// generateCreateMethods
func (r *GenerationRepo) generateCreateMethods() (string, error) {
	var createMethods string
	interfaceCreateOne, err := NewTemplate("InterfaceCreateOne").Parse(InterfaceCreateOne).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceCreateOne.String())
	interfaceCreateOneByTx, err := NewTemplate("InterfaceCreateOneByTx").Parse(InterfaceCreateOneByTx).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceCreateOneByTx.String())
	interfaceCreateBatch, err := NewTemplate("InterfaceCreateBatch").Parse(InterfaceCreateBatch).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceCreateBatch.String())
	return createMethods, nil
}

// generateUpdateMethods
func (r *GenerationRepo) generateUpdateMethods() (string, error) {
	var updateMethods string
	var primaryKey string
	for _, index := range r.index {
		isPrimaryKey, _ := index.PrimaryKey()
		if isPrimaryKey {
			primaryKey = index.Name()
			break
		}

	}
	if primaryKey == "" {
		return "", nil
	}
	interfaceUpdateOne, err := NewTemplate("InterfaceUpdateOne").Parse(InterfaceUpdateOne).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	updateMethods += fmt.Sprintln(interfaceUpdateOne.String())
	interfaceUpdateOneByTx, err := NewTemplate("InterfaceUpdateOneByTx").Parse(InterfaceUpdateOneByTx).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	updateMethods += fmt.Sprintln(interfaceUpdateOneByTx.String())
	return updateMethods, nil
}

// generateReadMethods
func (r *GenerationRepo) generateReadMethods() (string, error) {
	var readMethods string
	for _, v := range r.index {
		if r.CheckDaoFieldType(v.Columns()) {
			continue
		}
		unique, _ := v.Unique()
		// 唯一 && 字段数于1
		if unique && len(v.Columns()) == 1 {
			columnNameToDataType := r.columnNameToDataType[v.Columns()[0]]
			interfaceFindOneCacheByField, err := NewTemplate("InterfaceFindOneCacheByField").Parse(InterfaceFindOneCacheByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperFieldName(v.Columns()[0]),
				"lowerField":     r.LowerFieldName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindOneCacheByField.String())
			interfaceFindOneByField, err := NewTemplate("InterfaceFindOneByField").Parse(InterfaceFindOneByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperFieldName(v.Columns()[0]),
				"lowerField":     r.LowerFieldName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindOneByField.String())
			switch columnNameToDataType {
			case "bool":
			default:
				interfaceFindMultiCacheByFieldPlural, err := NewTemplate("InterfaceFindMultiCacheByFieldPlural").Parse(InterfaceFindMultiCacheByFieldPlural).Execute(map[string]any{
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				readMethods += fmt.Sprintln(interfaceFindMultiCacheByFieldPlural.String())
				interfaceFindMultiByFieldPlural, err := NewTemplate("InterfaceFindMultiByFieldPlural").Parse(InterfaceFindMultiByFieldPlural).Execute(map[string]any{
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				readMethods += fmt.Sprintln(interfaceFindMultiByFieldPlural.String())
			}

		}
		// 唯一 && 字段数大于1
		if unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			for _, vv := range v.Columns() {
				upperFields += r.UpperFieldName(vv)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerFieldName(vv), r.columnNameToDataType[vv])
			}
			interfaceFindOneCacheByFields, err := NewTemplate("InterfaceFindOneCacheByFields").Parse(InterfaceFindOneCacheByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			readMethods += fmt.Sprintln(interfaceFindOneCacheByFields.String())
			if err != nil {
				return "", err
			}
			interfaceFindOneByFields, err := NewTemplate("InterfaceFindOneByFields").Parse(InterfaceFindOneByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindOneByFields.String())
		}
		// 不唯一 && 字段数等于1
		if !unique && len(v.Columns()) == 1 {
			columnNameToDataType := r.columnNameToDataType[v.Columns()[0]]
			switch columnNameToDataType {
			case "bool":
			default:
				interfaceFindMultiByField, err := NewTemplate("InterfaceFindMultiByField").Parse(InterfaceFindMultiByField).Execute(map[string]any{
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				readMethods += fmt.Sprintln(interfaceFindMultiByField.String())
				interfaceFindMultiByFieldPlural, err := NewTemplate("InterfaceFindMultiByFieldPlural").Parse(InterfaceFindMultiByFieldPlural).Execute(map[string]any{
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				readMethods += fmt.Sprintln(interfaceFindMultiByFieldPlural.String())
			}
		}
		// 不唯一 && 字段数大于1
		if !unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			for _, v := range v.Columns() {
				upperFields += r.UpperFieldName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerFieldName(v), r.columnNameToDataType[v])
			}
			interfaceFindMultiByFields, err := NewTemplate("InterfaceFindMultiByFields").Parse(InterfaceFindMultiByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindMultiByFields.String())
		}
	}
	interfaceFindMultiByPaginator, err := NewTemplate("InterfaceFindMultiByPaginator").Parse(InterfaceFindMultiByPaginator).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	readMethods += fmt.Sprintln(interfaceFindMultiByPaginator.String())
	return readMethods, nil
}

// generateDelMethods
func (r *GenerationRepo) generateDelMethods() (string, error) {
	var delMethods string
	var haveUnique bool
	for _, v := range r.index {
		if r.CheckDaoFieldType(v.Columns()) {
			continue
		}
		unique, _ := v.Unique()
		if unique {
			haveUnique = true
		}
		// 唯一 && 字段数于1
		if unique && len(v.Columns()) == 1 {
			switch r.columnNameToDataType[v.Columns()[0]] {
			case "bool":
			default:
				interfaceDeleteOneCacheByField, err := NewTemplate("InterfaceDeleteOneCacheByField").Parse(InterfaceDeleteOneCacheByField).Execute(map[string]any{
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneCacheByField.String())
				interfaceDeleteOneCacheByFieldTx, err := NewTemplate("InterfaceDeleteOneCacheByFieldTx").Parse(InterfaceDeleteOneCacheByFieldTx).Execute(map[string]any{
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneCacheByFieldTx.String())
				interfaceDeleteOneByField, err := NewTemplate("InterfaceDeleteOneByField").Parse(InterfaceDeleteOneByField).Execute(map[string]any{
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneByField.String())
				interfaceDeleteOneByFieldTx, err := NewTemplate("InterfaceDeleteOneByFieldTx").Parse(InterfaceDeleteOneByFieldTx).Execute(map[string]any{
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneByFieldTx.String())
				interfaceDeleteMultiCacheByFieldPlural, err := NewTemplate("InterfaceDeleteMultiCacheByFieldPlural").Parse(InterfaceDeleteMultiCacheByFieldPlural).Execute(map[string]any{
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiCacheByFieldPlural.String())
				interfaceDeleteMultiCacheByFieldPluralTx, err := NewTemplate("InterfaceDeleteMultiCacheByFieldPluralTx").Parse(InterfaceDeleteMultiCacheByFieldPluralTx).Execute(map[string]any{
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiCacheByFieldPluralTx.String())
				interfaceDeleteMultiByFieldPlural, err := NewTemplate("InterfaceDeleteMultiByFieldPlural").Parse(InterfaceDeleteMultiByFieldPlural).Execute(map[string]any{
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiByFieldPlural.String())
				interfaceDeleteMultiByFieldPluralTx, err := NewTemplate("InterfaceDeleteMultiByFieldPluralTx").Parse(InterfaceDeleteMultiByFieldPluralTx).Execute(map[string]any{
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiByFieldPluralTx.String())
			}
		}
		// 唯一 && 字段数大于1
		if unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			for _, vv := range v.Columns() {
				upperFields += r.UpperFieldName(vv)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerFieldName(vv), r.columnNameToDataType[vv])
			}
			interfaceDeleteOneCacheByFields, err := NewTemplate("InterfaceDeleteOneCacheByFields").Parse(InterfaceDeleteOneCacheByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperFieldName(v.Columns()[0]),
				"lowerField":        r.LowerFieldName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneCacheByFields.String())
			interfaceDeleteOneCacheByFieldsTx, err := NewTemplate("InterfaceDeleteOneCacheByFieldsTx").Parse(InterfaceDeleteOneCacheByFieldsTx).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperFieldName(v.Columns()[0]),
				"lowerField":        r.LowerFieldName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneCacheByFieldsTx.String())
			interfaceDeleteOneByFields, err := NewTemplate("InterfaceDeleteOneByFields").Parse(InterfaceDeleteOneByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperFieldName(v.Columns()[0]),
				"lowerField":        r.LowerFieldName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneByFields.String())
			interfaceDeleteOneByFieldsTx, err := NewTemplate("InterfaceDeleteOneByFieldsTx").Parse(InterfaceDeleteOneByFieldsTx).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperFieldName(v.Columns()[0]),
				"lowerField":        r.LowerFieldName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneByFieldsTx.String())
		}
		// 不唯一 && 字段数等于1
		if !unique && len(v.Columns()) == 1 {

		}
		// 不唯一 && 字段数大于1
		if !unique && len(v.Columns()) > 1 {

		}
	}
	if !haveUnique {
		return "", nil
	}
	interfaceDeleteUniqueIndexCacheTpl, err := NewTemplate("InterfaceDeleteUniqueIndexCache").Parse(InterfaceDeleteUniqueIndexCache).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	delMethods += fmt.Sprintln(interfaceDeleteUniqueIndexCacheTpl.String())
	return delMethods, nil
}

// generateTypes
func (r *GenerationRepo) generateTypes() (string, error) {
	var methods string
	createMethods, err := r.generateCreateMethods()
	if err != nil {
		return "", err
	}
	updateMethods, err := r.generateUpdateMethods()
	if err != nil {
		return "", err
	}
	readMethods, err := r.generateReadMethods()
	if err != nil {
		return "", err
	}
	delMethods, err := r.generateDelMethods()
	if err != nil {
		return "", err
	}
	methods += createMethods
	methods += updateMethods
	methods += readMethods
	methods += delMethods
	typesTpl, err := NewTemplate("Types").Parse(Types).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
		"methods":        methods,
	})
	return typesTpl.String(), nil
}

// generateNew
func (r *GenerationRepo) generateNew() (string, error) {
	newTpl, err := NewTemplate("New").Parse(New).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	return newTpl.String(), nil
}

// generateCreateFunc
func (r *GenerationRepo) generateCreateFunc() (string, error) {
	var createFunc string
	createOne, err := NewTemplate("CreateOne").Parse(CreateOne).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(createOne.String())
	createOneByTx, err := NewTemplate("CreateOneByTx").Parse(CreateOneByTx).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(createOneByTx.String())
	createBatch, err := NewTemplate("CreateBatch").Parse(CreateBatch).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(createBatch.String())
	return createFunc, nil
}

// generateReadFunc
func (r *GenerationRepo) generateReadFunc() (string, error) {
	var readFunc string
	for _, v := range r.index {
		if r.CheckDaoFieldType(v.Columns()) {
			continue
		}
		unique, _ := v.Unique()
		// 唯一 && 字段数于1
		if unique && len(v.Columns()) == 1 {
			var whereField string
			columnNameToDataType := r.columnNameToDataType[v.Columns()[0]]
			switch columnNameToDataType {
			case "bool":
				whereField += fmt.Sprintf("dao.%s.Is(%s),", r.UpperFieldName(v.Columns()[0]), r.LowerFieldName(v.Columns()[0]))
			default:
				whereField += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperFieldName(v.Columns()[0]), r.LowerFieldName(v.Columns()[0]))
			}
			findOneCacheByField, err := NewTemplate("findOneCacheByField").Parse(FindOneCacheByField).Execute(map[string]any{
				"firstTableChar": r.firstTableChar,
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperFieldName(v.Columns()[0]),
				"lowerField":     r.LowerFieldName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
				"whereField":     whereField,
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findOneCacheByField.String())
			findOneByField, err := NewTemplate("findOneByField").Parse(FindOneByField).Execute(map[string]any{
				"firstTableChar": r.firstTableChar,
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperFieldName(v.Columns()[0]),
				"lowerField":     r.LowerFieldName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
				"whereField":     whereField,
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findOneByField.String())

			switch columnNameToDataType {
			case "bool":
			default:
				findMultiCacheByFieldPlural, err := NewTemplate("findMultiCacheByFieldPlural").Parse(FindMultiCacheByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
					"whereField":       whereField,
				})
				if err != nil {
					return "", err
				}
				readFunc += fmt.Sprintln(findMultiCacheByFieldPlural.String())
				findMultiByFieldPlural, err := NewTemplate("findMultiByFieldPlural").Parse(FindMultiByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"lowerDBName":      r.lowerDBName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
					"whereField":       whereField,
				})
				if err != nil {
					return "", err
				}
				readFunc += fmt.Sprintln(findMultiByFieldPlural.String())
			}
		}
		// 唯一 && 字段数大于1
		if unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var lowerFieldsJoin string
			var whereFields string
			for _, v := range v.Columns() {
				upperFields += r.UpperFieldName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerFieldName(v), r.columnNameToDataType[v])
				lowerFieldsJoin += fmt.Sprintf("%s,", r.LowerFieldName(v))
				switch r.columnNameToDataType[v] {
				case "bool":
					whereFields += fmt.Sprintf("dao.%s.Is(%s),", r.UpperFieldName(v), r.LowerFieldName(v))
				default:
					whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperFieldName(v), r.LowerFieldName(v))
				}
			}
			findOneCacheByFields, err := NewTemplate("findOneCacheByFields").Parse(FindOneCacheByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
				"lowerFieldsJoin":   strings.Trim(lowerFieldsJoin, ","),
			})
			readFunc += fmt.Sprintln(findOneCacheByFields.String())
			if err != nil {
				return "", err
			}
			findOneByFields, err := NewTemplate("findOneByFields").Parse(FindOneByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findOneByFields.String())
		}
		// 不唯一 && 字段数等于1
		if !unique && len(v.Columns()) == 1 {
			var whereField string
			columnNameToDataType := r.columnNameToDataType[v.Columns()[0]]
			switch columnNameToDataType {
			case "bool":
				whereField += fmt.Sprintf("dao.%s.Is(%s),", r.UpperFieldName(v.Columns()[0]), r.LowerFieldName(v.Columns()[0]))
			default:
				whereField += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperFieldName(v.Columns()[0]), r.LowerFieldName(v.Columns()[0]))
			}
			findMultiByField, err := NewTemplate("findMultiByField").Parse(FindMultiByField).Execute(map[string]any{
				"firstTableChar": r.firstTableChar,
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperFieldName(v.Columns()[0]),
				"lowerField":     r.LowerFieldName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
				"whereField":     whereField,
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findMultiByField.String())
			switch columnNameToDataType {
			case "bool":
			default:
				findMultiByFieldPlural, err := NewTemplate("findMultiByFieldPlural").Parse(FindMultiByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"lowerDBName":      r.lowerDBName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
					"whereField":       whereField,
				})
				if err != nil {
					return "", err
				}
				readFunc += fmt.Sprintln(findMultiByFieldPlural.String())
			}
		}
		// 不唯一 && 字段数大于1
		if !unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var whereFields string
			for _, v := range v.Columns() {
				upperFields += r.UpperFieldName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerFieldName(v), r.columnNameToDataType[v])
				switch r.columnNameToDataType[v] {
				case "bool":
					whereFields += fmt.Sprintf("dao.%s.Is(%s),", r.UpperFieldName(v), r.LowerFieldName(v))
				default:
					whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperFieldName(v), r.LowerFieldName(v))
				}
			}
			findMultiByFields, err := NewTemplate("findMultiByFields").Parse(FindMultiByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findMultiByFields.String())
		}
	}
	findMultiByPaginator, err := NewTemplate("FindMultiByPaginator").Parse(FindMultiByPaginator).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	readFunc += fmt.Sprintln(findMultiByPaginator.String())
	return readFunc, nil
}

// generateUpdateFunc
func (r *GenerationRepo) generateUpdateFunc() (string, error) {
	var updateFunc string
	var primaryKey string
	for _, index := range r.index {
		isPrimaryKey, _ := index.PrimaryKey()
		if isPrimaryKey {
			primaryKey = index.Columns()[0]
			break
		}
	}
	if primaryKey == "" {
		return "", nil
	}
	updateOneTpl, err := NewTemplate("UpdateOne").Parse(UpdateOne).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
		"upperField":     r.UpperFieldName(primaryKey),
	})
	if err != nil {
		return "", err
	}
	updateFunc += fmt.Sprintln(updateOneTpl.String())
	updateOneByTx, err := NewTemplate("UpdateOneByTx").Parse(UpdateOneByTx).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
		"upperField":     r.UpperFieldName(primaryKey),
	})
	if err != nil {
		return "", err
	}
	updateFunc += fmt.Sprintln(updateOneByTx.String())
	return updateFunc, nil
}

// generateDelFunc
func (r *GenerationRepo) generateDelFunc() (string, error) {
	var delMethods string
	var varCacheDelKeys string
	var haveUnique bool
	for _, v := range r.index {
		if r.CheckDaoFieldType(v.Columns()) {
			continue
		}
		unique, _ := v.Unique()
		if unique {
			haveUnique = true
			var cacheField string
			cacheFieldsJoinSli := make([]string, 0)
			for _, column := range v.Columns() {
				cacheField += r.UpperFieldName(column)
				cacheFieldsJoinSli = append(cacheFieldsJoinSli, fmt.Sprintf("v.%s", r.UpperFieldName(column)))
			}
			varCacheDelKeyTpl, err := NewTemplate("VarCacheDelKey").Parse(VarCacheDelKey).Execute(map[string]any{
				"firstTableChar":  r.firstTableChar,
				"upperTableName":  r.upperTableName,
				"cacheField":      cacheField,
				"cacheFieldsJoin": strings.Join(cacheFieldsJoinSli, ","),
			})
			if err != nil {
				return "", err
			}
			varCacheDelKeys += fmt.Sprintln(varCacheDelKeyTpl.String())
		}
		// 唯一 && 字段数于1
		if unique && len(v.Columns()) == 1 {
			columnNameToDataType := r.columnNameToDataType[v.Columns()[0]]
			switch columnNameToDataType {
			case "bool":
			default:
				deleteOneCacheByField, err := NewTemplate("DeleteOneCacheByField").Parse(DeleteOneCacheByField).Execute(map[string]any{
					"firstTableChar": r.firstTableChar,
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteOneCacheByField.String())
				deleteOneCacheByFieldTx, err := NewTemplate("DeleteOneCacheByFieldTx").Parse(DeleteOneCacheByFieldTx).Execute(map[string]any{
					"firstTableChar": r.firstTableChar,
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteOneCacheByFieldTx.String())
				deleteOneByField, err := NewTemplate("DeleteOneByField").Parse(DeleteOneByField).Execute(map[string]any{
					"firstTableChar": r.firstTableChar,
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteOneByField.String())
				deleteOneByFieldTx, err := NewTemplate("DeleteOneByFieldTx").Parse(DeleteOneByFieldTx).Execute(map[string]any{
					"firstTableChar": r.firstTableChar,
					"lowerDBName":    r.lowerDBName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.UpperFieldName(v.Columns()[0]),
					"lowerField":     r.LowerFieldName(v.Columns()[0]),
					"dataType":       r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteOneByFieldTx.String())
				deleteMultiCacheByFieldPlural, err := NewTemplate("DeleteMultiCacheByFieldPlural").Parse(DeleteMultiCacheByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteMultiCacheByFieldPlural.String())
				deleteMultiCacheByFieldPluralTx, err := NewTemplate("DeleteMultiCacheByFieldPluralTx").Parse(DeleteMultiCacheByFieldPluralTx).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteMultiCacheByFieldPluralTx.String())
				deleteMultiByFieldPlural, err := NewTemplate("DeleteMultiByFieldPlural").Parse(DeleteMultiByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteMultiByFieldPlural.String())
				deleteMultiByFieldPluralTx, err := NewTemplate("DeleteMultiByFieldPluralTx").Parse(DeleteMultiByFieldPluralTx).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"lowerDBName":      r.lowerDBName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.UpperFieldName(v.Columns()[0]),
					"lowerField":       r.LowerFieldName(v.Columns()[0]),
					"upperFieldPlural": r.Plural(r.UpperFieldName(v.Columns()[0])),
					"lowerFieldPlural": r.Plural(r.LowerFieldName(v.Columns()[0])),
					"dataType":         r.columnNameToDataType[v.Columns()[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteMultiByFieldPluralTx.String())
			}
		}
		// 唯一 && 字段数大于1
		if unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var whereFields string
			for _, v := range v.Columns() {
				upperFields += r.UpperFieldName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerFieldName(v), r.columnNameToDataType[v])
				switch r.columnNameToDataType[v] {
				case "bool":
					whereFields += fmt.Sprintf("dao.%s.Is(%s),", r.UpperFieldName(v), r.LowerFieldName(v))
				default:
					whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperFieldName(v), r.LowerFieldName(v))
				}
			}
			deleteOneCacheByFields, err := NewTemplate("DeleteOneCacheByFields").Parse(DeleteOneCacheByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperFieldName(v.Columns()[0]),
				"lowerField":        r.LowerFieldName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneCacheByFields.String())
			deleteOneCacheByFieldsTx, err := NewTemplate("DeleteOneCacheByFields").Parse(DeleteOneCacheByFieldsTx).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperFieldName(v.Columns()[0]),
				"lowerField":        r.LowerFieldName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneCacheByFieldsTx.String())
			deleteOneByFields, err := NewTemplate("DeleteOneByFields").Parse(DeleteOneByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"lowerField":        r.LowerFieldName(v.Columns()[0]),
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneByFields.String())
			deleteOneByFieldsTx, err := NewTemplate("DeleteOneByFieldsTx").Parse(DeleteOneByFieldsTx).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"lowerField":        r.LowerFieldName(v.Columns()[0]),
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneByFieldsTx.String())
		}
		// 不唯一 && 字段数等于1
		if !unique && len(v.Columns()) == 1 {

		}
		// 不唯一 && 字段数大于1
		if !unique && len(v.Columns()) > 1 {

		}
	}
	if !haveUnique {
		return "", nil
	}
	deleteUniqueIndexCacheTpl, err := NewTemplate("DeleteUniqueIndexCache").Parse(DeleteUniqueIndexCache).Execute(map[string]any{
		"firstTableChar":  r.firstTableChar,
		"lowerDBName":     r.lowerDBName,
		"upperTableName":  r.upperTableName,
		"lowerTableName":  r.lowerTableName,
		"varCacheDelKeys": varCacheDelKeys,
	})
	if err != nil {
		return "", err
	}
	delMethods += fmt.Sprintln(deleteUniqueIndexCacheTpl.String())
	return delMethods, nil
}

// UpperFieldName 字段名称大写
func (r *GenerationRepo) UpperFieldName(s string) string {
	return r.columnNameToName[s]
}

// LowerFieldName 字段名称小写
func (r *GenerationRepo) LowerFieldName(s string) string {
	str := r.UpperFieldName(s)
	if str == "" {
		return str
	}
	words := []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "ttl", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	// 如果第一个单词命中  则不处理
	for _, v := range words {
		if strings.HasPrefix(str, v) {
			return str
		}
	}
	rs := []rune(str)
	f := rs[0]
	if 'A' <= f && f <= 'Z' {
		str = string(unicode.ToLower(f)) + string(rs[1:])
	}
	if token.Lookup(str).IsKeyword() || StrSliFind(KeyWords, str) {
		str = "_" + str
	}
	return str
}

// UpperName 大写
func (r *GenerationRepo) UpperName(s string) string {
	return r.gorm.NamingStrategy.SchemaName(s)
}

// LowerName 小写
func (r *GenerationRepo) LowerName(s string) string {
	str := r.UpperName(s)
	if str == "" {
		return str
	}
	words := []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "ttl", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	// 如果第一个单词命中  则不处理
	for _, v := range words {
		if strings.HasPrefix(str, v) {
			return str
		}
	}
	rs := []rune(str)
	f := rs[0]
	if 'A' <= f && f <= 'Z' {
		str = string(unicode.ToLower(f)) + string(rs[1:])
	}
	return str
}

// Plural 复数形式
func (r *GenerationRepo) Plural(s string) string {
	str := inflection.Plural(s)
	if str == s {
		str += "Plural"
	}
	return str
}

// CheckDaoFieldType  检查字段状态
func (r *GenerationRepo) CheckDaoFieldType(s []string) bool {
	for _, v := range s {
		if r.columnNameToFieldType[v] == "Field" {
			return true
		}
	}
	return false
}
