//nolint:all
package repo

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
	"golang.org/x/tools/imports"
	"gorm.io/gorm"
)

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

func (r *Repo) GenerationTable(table string, columnNameToDataType map[string]string) error {
	var file string
	// 查询当前db的索引
	indexes, err := r.gorm.Migrator().GetIndexes(table)
	if err != nil {
		return err
	}
	// 索引按名称排序
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i].Name() > indexes[j].Name()
	})
	generationRepo := GenerationRepo{
		gorm:                 r.gorm,
		columnNameToDataType: columnNameToDataType,
		lowerDBName:          r.gorm.Migrator().CurrentDatabase(),
		lowerTableName:       "",
		upperTableName:       "",
		daoPkgPath:           FillModelPkgPath(r.daoPath),
		modelPkgPath:         FillModelPkgPath(r.modelPath),
		index:                indexes,
	}
	generationRepo.lowerTableName = generationRepo.LowerName(table)
	generationRepo.upperTableName = generationRepo.UpperName(table)
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
	gorm                 *gorm.DB
	columnNameToDataType map[string]string // 字段名称对应的类型
	lowerDBName          string
	lowerTableName       string
	upperTableName       string
	daoPkgPath           string
	modelPkgPath         string
	index                []gorm.Index
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
	var cacheKeys string
	for _, v := range r.index {
		unique, _ := v.Unique()
		if unique {
			var cacheField string
			for _, column := range v.Columns() {
				cacheField += r.UpperName(column)
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
	tplParams := map[string]any{
		"upperTableName": r.upperTableName,
		"cacheKeys":      cacheKeys,
	}
	tpl, err := NewTemplate("Var").Parse(Var).Execute(tplParams)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
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
	tpl, err := NewTemplate("InterfaceUpdateOne").Parse(InterfaceUpdateOne).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintln(tpl.String()), nil
}

// generateReadMethods
func (r *GenerationRepo) generateReadMethods() (string, error) {
	var readMethods string
	for _, v := range r.index {
		unique, _ := v.Unique()
		// 唯一 && 字段数于1
		if unique && len(v.Columns()) == 1 {
			interfaceFindOneCacheByField, err := NewTemplate("InterfaceFindOneCacheByField").Parse(InterfaceFindOneCacheByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
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
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindOneByField.String())
			interfaceFindMultiCacheByFieldPlural, err := NewTemplate("InterfaceFindMultiCacheByFieldPlural").Parse(InterfaceFindMultiCacheByFieldPlural).Execute(map[string]any{
				"lowerDBName":      r.lowerDBName,
				"upperTableName":   r.upperTableName,
				"lowerTableName":   r.lowerTableName,
				"upperField":       r.UpperName(v.Columns()[0]),
				"lowerField":       r.LowerName(v.Columns()[0]),
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
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
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindMultiByFieldPlural.String())
		}
		// 唯一 && 字段数大于1
		if unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			for _, v := range v.Columns() {
				upperFields += r.UpperName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerName(v), r.columnNameToDataType[v])
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
			interfaceFindMultiByField, err := NewTemplate("InterfaceFindMultiByField").Parse(InterfaceFindMultiByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
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
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindMultiByFieldPlural.String())
		}
		// 不唯一 && 字段数大于1
		if !unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			for _, v := range v.Columns() {
				upperFields += r.UpperName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerName(v), r.columnNameToDataType[v])
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
	return readMethods, nil
}

// generateDelMethods
func (r *GenerationRepo) generateDelMethods() (string, error) {
	var delMethods string
	for _, v := range r.index {
		unique, _ := v.Unique()
		// 唯一 && 字段数于1
		if unique && len(v.Columns()) == 1 {
			interfaceDeleteOneCacheByField, err := NewTemplate("InterfaceDeleteOneCacheByField").Parse(InterfaceDeleteOneCacheByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneCacheByField.String())
			interfaceDeleteOneByField, err := NewTemplate("InterfaceDeleteOneByField").Parse(InterfaceDeleteOneByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneByField.String())

			interfaceDeleteMultiCacheByFieldPlural, err := NewTemplate("InterfaceDeleteMultiCacheByFieldPlural").Parse(InterfaceDeleteMultiCacheByFieldPlural).Execute(map[string]any{
				"lowerDBName":      r.lowerDBName,
				"upperTableName":   r.upperTableName,
				"lowerTableName":   r.lowerTableName,
				"upperField":       r.UpperName(v.Columns()[0]),
				"lowerField":       r.LowerName(v.Columns()[0]),
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteMultiCacheByFieldPlural.String())
			interfaceDeleteMultiByFieldPlural, err := NewTemplate("InterfaceDeleteMultiByFieldPlural").Parse(InterfaceDeleteMultiByFieldPlural).Execute(map[string]any{
				"lowerDBName":      r.lowerDBName,
				"upperTableName":   r.upperTableName,
				"lowerTableName":   r.lowerTableName,
				"upperField":       r.UpperName(v.Columns()[0]),
				"lowerField":       r.LowerName(v.Columns()[0]),
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteMultiByFieldPlural.String())
		}
		// 唯一 && 字段数大于1
		if unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var whereFields string
			for _, v := range v.Columns() {
				upperFields += r.UpperName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerName(v), r.columnNameToDataType[v])
				whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperName(v), r.LowerName(v))
			}
			interfaceDeleteOneCacheByFields, err := NewTemplate("InterfaceDeleteOneCacheByFields").Parse(InterfaceDeleteOneCacheByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperName(v.Columns()[0]),
				"lowerField":        r.LowerName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneCacheByFields.String())
			interfaceDeleteOneByFields, err := NewTemplate("InterfaceDeleteOneByFields").Parse(InterfaceDeleteOneByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperName(v.Columns()[0]),
				"lowerField":        r.LowerName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneByFields.String())
		}
		// 不唯一 && 字段数等于1
		if !unique && len(v.Columns()) == 1 {

		}
		// 不唯一 && 字段数大于1
		if !unique && len(v.Columns()) > 1 {

		}
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
	createOneTpl, err := NewTemplate("CreateOne").Parse(CreateOne).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(createOneTpl.String())
	createBatchTpl, err := NewTemplate("CreateBatch").Parse(CreateBatch).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(createBatchTpl.String())
	return createFunc, nil
}

// generateReadFunc
func (r *GenerationRepo) generateReadFunc() (string, error) {
	var readFunc string
	for _, v := range r.index {
		unique, _ := v.Unique()
		// 唯一 && 字段数于1
		if unique && len(v.Columns()) == 1 {
			findOneCacheByField, err := NewTemplate("findOneCacheByField").Parse(FindOneCacheByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findOneCacheByField.String())
			findOneByField, err := NewTemplate("findOneByField").Parse(FindOneByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findOneByField.String())
			findMultiCacheByFieldPlural, err := NewTemplate("findMultiCacheByFieldPlural").Parse(FindMultiCacheByFieldPlural).Execute(map[string]any{
				"lowerDBName":      r.lowerDBName,
				"upperTableName":   r.upperTableName,
				"lowerTableName":   r.lowerTableName,
				"upperField":       r.UpperName(v.Columns()[0]),
				"lowerField":       r.LowerName(v.Columns()[0]),
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findMultiCacheByFieldPlural.String())

			findMultiByFieldPlural, err := NewTemplate("findMultiByFieldPlural").Parse(FindMultiByFieldPlural).Execute(map[string]any{
				"lowerDBName":      r.lowerDBName,
				"upperField":       r.UpperName(v.Columns()[0]),
				"upperTableName":   r.upperTableName,
				"lowerTableName":   r.lowerTableName,
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findMultiByFieldPlural.String())

		}
		// 唯一 && 字段数大于1
		if unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var lowerFieldsJoin string
			var whereFields string
			for _, v := range v.Columns() {
				upperFields += r.UpperName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerName(v), r.columnNameToDataType[v])
				lowerFieldsJoin += fmt.Sprintf("%s,", r.LowerName(v))
				whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperName(v), r.LowerName(v))
			}
			findOneCacheByFields, err := NewTemplate("findOneCacheByFields").Parse(FindOneCacheByFields).Execute(map[string]any{
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
			findMultiByField, err := NewTemplate("findMultiByField").Parse(FindMultiByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findMultiByField.String())
			findMultiByFieldPlural, err := NewTemplate("findMultiByFieldPlural").Parse(FindMultiByFieldPlural).Execute(map[string]any{
				"lowerDBName":      r.lowerDBName,
				"upperField":       r.UpperName(v.Columns()[0]),
				"upperTableName":   r.upperTableName,
				"lowerTableName":   r.lowerTableName,
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findMultiByFieldPlural.String())
		}
		// 不唯一 && 字段数大于1
		if !unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var whereFields string
			for _, v := range v.Columns() {
				upperFields += r.UpperName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerName(v), r.columnNameToDataType[v])
				whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperName(v), r.LowerName(v))
			}
			findMultiByFields, err := NewTemplate("findMultiByFields").Parse(FindMultiByFields).Execute(map[string]any{
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
	return readFunc, nil
}

// generateUpdateFunc
func (r *GenerationRepo) generateUpdateFunc() (string, error) {
	updateOneTpl, err := NewTemplate("UpdateOne").Parse(UpdateOne).Execute(map[string]any{
		"lowerDBName":    r.lowerDBName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintln(updateOneTpl.String()), nil
}

// generateDelFunc
func (r *GenerationRepo) generateDelFunc() (string, error) {
	var delMethods string
	var varCacheDelKeys string
	for _, v := range r.index {
		unique, _ := v.Unique()
		if unique {
			var cacheField string
			cacheFieldsJoinSli := make([]string, 0)
			for _, column := range v.Columns() {
				cacheField += r.UpperName(column)
				cacheFieldsJoinSli = append(cacheFieldsJoinSli, fmt.Sprintf("v.%s", r.UpperName(column)))
			}
			varCacheDelKeyTpl, err := NewTemplate("VarCacheDelKey").Parse(VarCacheDelKey).Execute(map[string]any{
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
			deleteOneCacheByField, err := NewTemplate("DeleteOneCacheByField").Parse(DeleteOneCacheByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneCacheByField.String())
			deleteOneByField, err := NewTemplate("DeleteOneByField").Parse(DeleteOneByField).Execute(map[string]any{
				"lowerDBName":    r.lowerDBName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.UpperName(v.Columns()[0]),
				"lowerField":     r.LowerName(v.Columns()[0]),
				"dataType":       r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneByField.String())
			deleteMultiCacheByFieldPlural, err := NewTemplate("DeleteMultiCacheByFieldPlural").Parse(DeleteMultiCacheByFieldPlural).Execute(map[string]any{
				"lowerDBName":      r.lowerDBName,
				"upperTableName":   r.upperTableName,
				"lowerTableName":   r.lowerTableName,
				"upperField":       r.UpperName(v.Columns()[0]),
				"lowerField":       r.LowerName(v.Columns()[0]),
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteMultiCacheByFieldPlural.String())
			deleteMultiByFieldPlural, err := NewTemplate("DeleteMultiByFieldPlural").Parse(DeleteMultiByFieldPlural).Execute(map[string]any{
				"lowerDBName":      r.lowerDBName,
				"upperTableName":   r.upperTableName,
				"lowerTableName":   r.lowerTableName,
				"upperField":       r.UpperName(v.Columns()[0]),
				"lowerField":       r.LowerName(v.Columns()[0]),
				"upperFieldPlural": inflection.Plural(r.UpperName(v.Columns()[0])),
				"lowerFieldPlural": inflection.Plural(r.LowerName(v.Columns()[0])),
				"dataType":         r.columnNameToDataType[v.Columns()[0]],
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteMultiByFieldPlural.String())
		}
		// 唯一 && 字段数大于1
		if unique && len(v.Columns()) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var whereFields string
			for _, v := range v.Columns() {
				upperFields += r.UpperName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.LowerName(v), r.columnNameToDataType[v])
				whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.UpperName(v), r.LowerName(v))
			}
			deleteOneCacheByFields, err := NewTemplate("DeleteOneCacheByFields").Parse(DeleteOneCacheByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.UpperName(v.Columns()[0]),
				"lowerField":        r.LowerName(v.Columns()[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns()[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneCacheByFields.String())
			deleteOneByFields, err := NewTemplate("DeleteOneByFields").Parse(DeleteOneByFields).Execute(map[string]any{
				"lowerDBName":       r.lowerDBName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"lowerField":        r.LowerName(v.Columns()[0]),
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneByFields.String())
		}
		// 不唯一 && 字段数等于1
		if !unique && len(v.Columns()) == 1 {

		}
		// 不唯一 && 字段数大于1
		if !unique && len(v.Columns()) > 1 {

		}
	}
	deleteUniqueIndexCacheTpl, err := NewTemplate("DeleteUniqueIndexCache").Parse(DeleteUniqueIndexCache).Execute(map[string]any{
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

// UpperName 大写
func (r *GenerationRepo) UpperName(s string) string {
	return r.gorm.NamingStrategy.SchemaName(s)
}

// LowerName 小写
func (r *GenerationRepo) LowerName(s string) string {
	s = r.UpperName(s)
	if s == "" {
		return s
	}
	words := []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "ttl", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	// 如果第一个单词命中  则不处理
	for _, v := range words {
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
