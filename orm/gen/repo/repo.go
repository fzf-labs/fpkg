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

	"github.com/fzf-labs/fpkg/orm/gen/utils/dbfunc"
	"github.com/fzf-labs/fpkg/orm/gen/utils/template"
	"github.com/fzf-labs/fpkg/orm/gen/utils/util"
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

func GenerationTable(db *gorm.DB, dbname, daoPath, modelPath, repoPath, table string, columnNameToDataType, columnNameToName, columnNameToFieldType map[string]string) error {
	var file string
	g := Repo{
		gorm:                  db,
		daoPath:               daoPath,
		modelPath:             modelPath,
		repoPath:              repoPath,
		table:                 table,
		columnNameToDataType:  columnNameToDataType,
		columnNameToName:      columnNameToName,
		columnNameToFieldType: columnNameToFieldType,
		firstTableChar:        "",
		dbName:                dbname,
		lowerTableName:        "",
		upperTableName:        "",
		daoPkgPath:            util.FillModelPkgPath(daoPath),
		modelPkgPath:          util.FillModelPkgPath(modelPath),
		index:                 make([]DBIndex, 0),
	}
	// 查询当前db的索引
	index, err := g.processIndex()
	if err != nil {
		return err
	}
	g.index = index
	g.lowerTableName = g.lowerName(table)
	g.upperTableName = g.upperName(table)
	g.firstTableChar = g.lowerTableName[0:1]
	generatePkg, err := g.generatePkg()
	if err != nil {
		return err
	}
	generateImport, err := g.generateImport()
	if err != nil {
		return err
	}
	generateVar, err := g.generateVar()
	if err != nil {
		return err
	}
	generateTypes, err := g.generateTypes()
	if err != nil {
		return err
	}
	generateNew, err := g.generateNew()
	if err != nil {
		return err
	}
	generateCreateFunc, err := g.generateCreateFunc()
	if err != nil {
		return err
	}
	generateUpdateFunc, err := g.generateUpdateFunc()
	if err != nil {
		return err
	}
	generateDelFunc, err := g.generateDelFunc()
	if err != nil {
		return err
	}
	generateReadFunc, err := g.generateReadFunc()
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
	outputFile := g.repoPath + "/" + table + ".repo.go"
	err = g.output(outputFile, []byte(file))
	if err != nil {
		return err
	}
	return nil
}

type Repo struct {
	gorm                  *gorm.DB          // 数据库
	daoPath               string            // dao所在的路径
	modelPath             string            // model所在的路径
	repoPath              string            // repo所在的路径
	table                 string            // 表名称
	columnNameToDataType  map[string]string // 字段名称对应的类型
	columnNameToName      map[string]string // 字段名称对应的Go名称
	columnNameToFieldType map[string]string // 字段名称对应的dao类型
	dbName                string            // 数据库名称
	firstTableChar        string            // 表名称第一个字母
	lowerTableName        string            // 表名称小写
	upperTableName        string            // 表名称大写
	daoPkgPath            string            // go文件中daoPkgPath
	modelPkgPath          string            // go文件中modelPkgPath
	index                 []DBIndex         // 索引
}

type DBIndex struct {
	Name       string   // 索引名称
	ColumnKey  string   // 索引字段KEY
	PrimaryKey bool     // 是否是主键
	Unique     bool     // 是否是唯一索引
	Columns    []string // 索引字段
}

// processIndex 索引处理  索引去重和排序
func (r *Repo) processIndex() ([]DBIndex, error) {
	result := make([]DBIndex, 0)
	tmp := make([]DBIndex, 0)
	repeat := make(map[string]struct{})
	// 查询是否有分区表
	childTableForTable, err := dbfunc.GetPartitionChildTableForTable(r.gorm, r.table)
	if err != nil {
		return nil, err
	}
	table := r.table
	if len(childTableForTable) > 0 {
		table = childTableForTable[0]
	}
	// 获取索引
	indexes, err := r.gorm.Migrator().GetIndexes(table)
	if err != nil {
		return nil, err
	}
	// 获取排序的索引字段
	sortIndexColumns, err := dbfunc.SortIndexColumns(r.gorm, table)
	if err != nil {
		return nil, err
	}
	for _, v := range indexes {
		primaryKey, _ := v.PrimaryKey()
		unique, _ := v.Unique()
		columns := sortIndexColumns[v.Name()]
		tmp = append(tmp, DBIndex{
			Name:       v.Name(),
			ColumnKey:  strings.Join(columns, "_"),
			PrimaryKey: primaryKey,
			Unique:     unique,
			Columns:    columns,
		})
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].ColumnKey < tmp[j].ColumnKey
	})
	// 主键索引
	for _, v := range tmp {
		if v.PrimaryKey {
			_, ok := repeat[v.ColumnKey]
			if !ok {
				repeat[v.ColumnKey] = struct{}{}
				result = append(result, v)
			}
		}
	}
	// 唯一索引
	for _, v := range tmp {
		if !v.PrimaryKey && v.Unique {
			_, ok := repeat[v.ColumnKey]
			if !ok {
				repeat[v.ColumnKey] = struct{}{}
				result = append(result, v)
			}

		}
	}
	// 普通索引
	for _, v := range tmp {
		if !v.PrimaryKey && !v.Unique {
			_, ok := repeat[v.ColumnKey]
			if !ok {
				repeat[v.ColumnKey] = struct{}{}
				result = append(result, v)
			}
		}
	}
	// 最左匹配原则索引
	for _, v := range tmp {
		if !v.PrimaryKey && len(v.Columns) > 1 {
			for i := len(v.Columns); i > 0; i-- {
				columnKey := strings.Join(v.Columns[0:i], "_")
				_, ok := repeat[columnKey]
				if !ok {
					repeat[columnKey] = struct{}{}
					result = append(result, DBIndex{
						Name:       v.Name,
						ColumnKey:  columnKey,
						PrimaryKey: false,
						Unique:     false,
						Columns:    v.Columns[0:i],
					})
				}
			}
		}
	}
	return result, nil
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

// generatePkg
func (r *Repo) generatePkg() (string, error) {
	tplParams := map[string]any{
		"dbName": r.dbName,
	}
	tpl, err := template.NewTemplate("Pkg").Parse(Pkg).Execute(tplParams)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

// generateImport
func (r *Repo) generateImport() (string, error) {
	tplParams := map[string]any{
		"daoPkgPath":   r.daoPkgPath,
		"modelPkgPath": r.modelPkgPath,
	}
	tpl, err := template.NewTemplate("Import").Parse(Import).Execute(tplParams)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

// generateVar
func (r *Repo) generateVar() (string, error) {
	var varStr string
	var cacheKeys string
	for _, v := range r.index {
		if r.checkDaoFieldType(v.Columns) {
			continue
		}
		if v.Unique {
			var cacheField string
			for _, column := range v.Columns {
				cacheField += r.upperFieldName(column)
			}
			varCacheTpl, err := template.NewTemplate("VarCache").Parse(VarCache).Execute(map[string]any{
				"dbName":         r.dbName,
				"upperTableName": r.upperTableName,
				"cacheField":     cacheField,
			})
			if err != nil {
				return "", err
			}
			cacheKeys += varCacheTpl.String()
		}
	}
	varTpl, err := template.NewTemplate("Var").Parse(Var).Execute(map[string]any{
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	varStr += fmt.Sprintln(varTpl.String())
	if len(cacheKeys) > 0 {
		varCacheKeysTpl, err := template.NewTemplate("Var").Parse(VarCacheKeys).Execute(map[string]any{
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
func (r *Repo) generateCreateMethods() (string, error) {
	var createMethods string
	interfaceCreateOne, err := template.NewTemplate("InterfaceCreateOne").Parse(InterfaceCreateOne).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceCreateOne.String())
	interfaceCreateOneByTx, err := template.NewTemplate("InterfaceCreateOneByTx").Parse(InterfaceCreateOneByTx).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceCreateOneByTx.String())
	interfaceUpsertOne, err := template.NewTemplate("InterfaceUpsertOne").Parse(InterfaceUpsertOne).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceUpsertOne.String())
	interfaceUpsertOneByTx, err := template.NewTemplate("InterfaceUpsertOneByTx").Parse(InterfaceUpsertOneByTx).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceUpsertOneByTx.String())
	interfaceUpsertOneByFields, err := template.NewTemplate("InterfaceUpsertOneByFields").Parse(InterfaceUpsertOneByFields).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceUpsertOneByFields.String())
	interfaceUpsertOneByFieldsTx, err := template.NewTemplate("InterfaceUpsertOneByFieldsTx").Parse(InterfaceUpsertOneByFieldsTx).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceUpsertOneByFieldsTx.String())
	interfaceCreateBatch, err := template.NewTemplate("InterfaceCreateBatch").Parse(InterfaceCreateBatch).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createMethods += fmt.Sprintln(interfaceCreateBatch.String())
	return createMethods, nil
}

// generateUpdateMethods
func (r *Repo) generateUpdateMethods() (string, error) {
	var updateMethods string
	var primaryKey string
	for _, v := range r.index {
		if v.PrimaryKey {
			primaryKey = v.Columns[0]
			break
		}
	}
	if primaryKey == "" {
		return "", nil
	}
	interfaceUpdateOne, err := template.NewTemplate("InterfaceUpdateOne").Parse(InterfaceUpdateOne).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	updateMethods += fmt.Sprintln(interfaceUpdateOne.String())
	interfaceUpdateOneByTx, err := template.NewTemplate("InterfaceUpdateOneByTx").Parse(InterfaceUpdateOneByTx).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	updateMethods += fmt.Sprintln(interfaceUpdateOneByTx.String())
	interfaceUpdateOneWithZero, err := template.NewTemplate("InterfaceUpdateOneWithZero").Parse(InterfaceUpdateOneWithZero).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	updateMethods += fmt.Sprintln(interfaceUpdateOneWithZero.String())
	interfaceUpdateOneWithZeroByTx, err := template.NewTemplate("InterfaceUpdateOneWithZeroByTx").Parse(InterfaceUpdateOneWithZeroByTx).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	updateMethods += fmt.Sprintln(interfaceUpdateOneWithZeroByTx.String())
	return updateMethods, nil
}

// generateReadMethods
func (r *Repo) generateReadMethods() (string, error) {
	var readMethods string
	for _, v := range r.index {
		if r.checkDaoFieldType(v.Columns) {
			continue
		}
		// 唯一 && 字段数于1
		if v.Unique && len(v.Columns) == 1 {
			columnNameToDataType := r.columnNameToDataType[v.Columns[0]]
			interfaceFindOneCacheByField, err := template.NewTemplate("InterfaceFindOneCacheByField").Parse(InterfaceFindOneCacheByField).Execute(map[string]any{
				"dbName":         r.dbName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.upperFieldName(v.Columns[0]),
				"lowerField":     r.lowerFieldName(v.Columns[0]),
				"dataType":       r.columnNameToDataType[v.Columns[0]],
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindOneCacheByField.String())
			interfaceFindOneByField, err := template.NewTemplate("InterfaceFindOneByField").Parse(InterfaceFindOneByField).Execute(map[string]any{
				"dbName":         r.dbName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.upperFieldName(v.Columns[0]),
				"lowerField":     r.lowerFieldName(v.Columns[0]),
				"dataType":       r.columnNameToDataType[v.Columns[0]],
			})
			if err != nil {
				return "", err
			}
			readMethods += fmt.Sprintln(interfaceFindOneByField.String())
			switch columnNameToDataType {
			case "bool":
			default:
				interfaceFindMultiCacheByFieldPlural, err := template.NewTemplate("InterfaceFindMultiCacheByFieldPlural").Parse(InterfaceFindMultiCacheByFieldPlural).Execute(map[string]any{
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				readMethods += fmt.Sprintln(interfaceFindMultiCacheByFieldPlural.String())
				interfaceFindMultiByFieldPlural, err := template.NewTemplate("InterfaceFindMultiByFieldPlural").Parse(InterfaceFindMultiByFieldPlural).Execute(map[string]any{
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				readMethods += fmt.Sprintln(interfaceFindMultiByFieldPlural.String())
			}

		}
		// 唯一 && 字段数大于1
		if v.Unique && len(v.Columns) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			for _, vv := range v.Columns {
				upperFields += r.upperFieldName(vv)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.lowerFieldName(vv), r.columnNameToDataType[vv])
			}
			interfaceFindOneCacheByFields, err := template.NewTemplate("InterfaceFindOneCacheByFields").Parse(InterfaceFindOneCacheByFields).Execute(map[string]any{
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			readMethods += fmt.Sprintln(interfaceFindOneCacheByFields.String())
			if err != nil {
				return "", err
			}
			interfaceFindOneByFields, err := template.NewTemplate("InterfaceFindOneByFields").Parse(InterfaceFindOneByFields).Execute(map[string]any{
				"dbName":            r.dbName,
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
		if !v.Unique && len(v.Columns) == 1 {
			columnNameToDataType := r.columnNameToDataType[v.Columns[0]]
			switch columnNameToDataType {
			case "bool":
			default:
				interfaceFindMultiByField, err := template.NewTemplate("InterfaceFindMultiByField").Parse(InterfaceFindMultiByField).Execute(map[string]any{
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				readMethods += fmt.Sprintln(interfaceFindMultiByField.String())
				interfaceFindMultiByFieldPlural, err := template.NewTemplate("InterfaceFindMultiByFieldPlural").Parse(InterfaceFindMultiByFieldPlural).Execute(map[string]any{
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				readMethods += fmt.Sprintln(interfaceFindMultiByFieldPlural.String())
			}
		}
		// 不唯一 && 字段数大于1
		if !v.Unique && len(v.Columns) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			for _, v := range v.Columns {
				upperFields += r.upperFieldName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.lowerFieldName(v), r.columnNameToDataType[v])
			}
			interfaceFindMultiByFields, err := template.NewTemplate("InterfaceFindMultiByFields").Parse(InterfaceFindMultiByFields).Execute(map[string]any{
				"dbName":            r.dbName,
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
	interfaceFindMultiByCustom, err := template.NewTemplate("InterfaceFindMultiByCustom").Parse(InterfaceFindMultiByCustom).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	readMethods += fmt.Sprintln(interfaceFindMultiByCustom.String())
	return readMethods, nil
}

// generateDelMethods
func (r *Repo) generateDelMethods() (string, error) {
	var delMethods string
	var haveUnique bool
	for _, v := range r.index {
		if r.checkDaoFieldType(v.Columns) {
			continue
		}
		if v.Unique {
			haveUnique = true
		}
		// 唯一 && 字段数于1
		if v.Unique && len(v.Columns) == 1 {
			switch r.columnNameToDataType[v.Columns[0]] {
			case "bool":
			default:
				interfaceDeleteOneCacheByField, err := template.NewTemplate("InterfaceDeleteOneCacheByField").Parse(InterfaceDeleteOneCacheByField).Execute(map[string]any{
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneCacheByField.String())
				interfaceDeleteOneCacheByFieldTx, err := template.NewTemplate("InterfaceDeleteOneCacheByFieldTx").Parse(InterfaceDeleteOneCacheByFieldTx).Execute(map[string]any{
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneCacheByFieldTx.String())
				interfaceDeleteOneByField, err := template.NewTemplate("InterfaceDeleteOneByField").Parse(InterfaceDeleteOneByField).Execute(map[string]any{
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneByField.String())
				interfaceDeleteOneByFieldTx, err := template.NewTemplate("InterfaceDeleteOneByFieldTx").Parse(InterfaceDeleteOneByFieldTx).Execute(map[string]any{
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteOneByFieldTx.String())
				interfaceDeleteMultiCacheByFieldPlural, err := template.NewTemplate("InterfaceDeleteMultiCacheByFieldPlural").Parse(InterfaceDeleteMultiCacheByFieldPlural).Execute(map[string]any{
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiCacheByFieldPlural.String())
				interfaceDeleteMultiCacheByFieldPluralTx, err := template.NewTemplate("InterfaceDeleteMultiCacheByFieldPluralTx").Parse(InterfaceDeleteMultiCacheByFieldPluralTx).Execute(map[string]any{
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiCacheByFieldPluralTx.String())
				interfaceDeleteMultiByFieldPlural, err := template.NewTemplate("InterfaceDeleteMultiByFieldPlural").Parse(InterfaceDeleteMultiByFieldPlural).Execute(map[string]any{
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiByFieldPlural.String())
				interfaceDeleteMultiByFieldPluralTx, err := template.NewTemplate("InterfaceDeleteMultiByFieldPluralTx").Parse(InterfaceDeleteMultiByFieldPluralTx).Execute(map[string]any{
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(interfaceDeleteMultiByFieldPluralTx.String())
			}
		}
		// 唯一 && 字段数大于1
		if v.Unique && len(v.Columns) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			for _, vv := range v.Columns {
				upperFields += r.upperFieldName(vv)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.lowerFieldName(vv), r.columnNameToDataType[vv])
			}
			interfaceDeleteOneCacheByFields, err := template.NewTemplate("InterfaceDeleteOneCacheByFields").Parse(InterfaceDeleteOneCacheByFields).Execute(map[string]any{
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.upperFieldName(v.Columns[0]),
				"lowerField":        r.lowerFieldName(v.Columns[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneCacheByFields.String())
			interfaceDeleteOneCacheByFieldsTx, err := template.NewTemplate("InterfaceDeleteOneCacheByFieldsTx").Parse(InterfaceDeleteOneCacheByFieldsTx).Execute(map[string]any{
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.upperFieldName(v.Columns[0]),
				"lowerField":        r.lowerFieldName(v.Columns[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneCacheByFieldsTx.String())
			interfaceDeleteOneByFields, err := template.NewTemplate("InterfaceDeleteOneByFields").Parse(InterfaceDeleteOneByFields).Execute(map[string]any{
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.upperFieldName(v.Columns[0]),
				"lowerField":        r.lowerFieldName(v.Columns[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneByFields.String())
			interfaceDeleteOneByFieldsTx, err := template.NewTemplate("InterfaceDeleteOneByFieldsTx").Parse(InterfaceDeleteOneByFieldsTx).Execute(map[string]any{
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.upperFieldName(v.Columns[0]),
				"lowerField":        r.lowerFieldName(v.Columns[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(interfaceDeleteOneByFieldsTx.String())
		}
		// 不唯一 && 字段数等于1
		if !v.Unique && len(v.Columns) == 1 {

		}
		// 不唯一 && 字段数大于1
		if !v.Unique && len(v.Columns) > 1 {

		}
	}
	if !haveUnique {
		return "", nil
	}
	interfaceDeleteUniqueIndexCacheTpl, err := template.NewTemplate("InterfaceDeleteUniqueIndexCache").Parse(InterfaceDeleteUniqueIndexCache).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	delMethods += fmt.Sprintln(interfaceDeleteUniqueIndexCacheTpl.String())
	return delMethods, nil
}

// generateTypes
func (r *Repo) generateTypes() (string, error) {
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
	typesTpl, err := template.NewTemplate("Types").Parse(Types).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
		"methods":        methods,
	})
	return typesTpl.String(), nil
}

// generateNew
func (r *Repo) generateNew() (string, error) {
	newTpl, err := template.NewTemplate("New").Parse(New).Execute(map[string]any{
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	return newTpl.String(), nil
}

// generateCreateFunc
func (r *Repo) generateCreateFunc() (string, error) {
	var createFunc string
	createOne, err := template.NewTemplate("CreateOne").Parse(CreateOne).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(createOne.String())
	createOneByTx, err := template.NewTemplate("CreateOneByTx").Parse(CreateOneByTx).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(createOneByTx.String())
	upsertOne, err := template.NewTemplate("UpsertOne").Parse(UpsertOne).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(upsertOne.String())
	upsertOneByTx, err := template.NewTemplate("UpsertOneByTx").Parse(UpsertOneByTx).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(upsertOneByTx.String())
	upsertOneByFields, err := template.NewTemplate("UpsertOneByFields").Parse(UpsertOneByFields).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(upsertOneByFields.String())
	upsertOneByFieldsTx, err := template.NewTemplate("UpsertOneByFieldsTx").Parse(UpsertOneByFieldsTx).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(upsertOneByFieldsTx.String())
	createBatch, err := template.NewTemplate("CreateBatch").Parse(CreateBatch).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
	})
	if err != nil {
		return "", err
	}
	createFunc += fmt.Sprintln(createBatch.String())
	return createFunc, nil
}

// generateReadFunc
func (r *Repo) generateReadFunc() (string, error) {
	var readFunc string
	for _, v := range r.index {
		if r.checkDaoFieldType(v.Columns) {
			continue
		}
		// 唯一 && 字段数于1
		if v.Unique && len(v.Columns) == 1 {
			var whereField string
			columnNameToDataType := r.columnNameToDataType[v.Columns[0]]
			switch columnNameToDataType {
			case "bool":
				whereField += fmt.Sprintf("dao.%s.Is(%s),", r.upperFieldName(v.Columns[0]), r.lowerFieldName(v.Columns[0]))
			default:
				whereField += fmt.Sprintf("dao.%s.Eq(%s),", r.upperFieldName(v.Columns[0]), r.lowerFieldName(v.Columns[0]))
			}
			findOneCacheByField, err := template.NewTemplate("findOneCacheByField").Parse(FindOneCacheByField).Execute(map[string]any{
				"firstTableChar": r.firstTableChar,
				"dbName":         r.dbName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.upperFieldName(v.Columns[0]),
				"lowerField":     r.lowerFieldName(v.Columns[0]),
				"dataType":       r.columnNameToDataType[v.Columns[0]],
				"whereField":     whereField,
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findOneCacheByField.String())
			findOneByField, err := template.NewTemplate("findOneByField").Parse(FindOneByField).Execute(map[string]any{
				"firstTableChar": r.firstTableChar,
				"dbName":         r.dbName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.upperFieldName(v.Columns[0]),
				"lowerField":     r.lowerFieldName(v.Columns[0]),
				"dataType":       r.columnNameToDataType[v.Columns[0]],
				"whereField":     whereField,
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findOneByField.String())

			switch columnNameToDataType {
			case "bool":
			default:
				findMultiCacheByFieldPlural, err := template.NewTemplate("findMultiCacheByFieldPlural").Parse(FindMultiCacheByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
					"whereField":       whereField,
				})
				if err != nil {
					return "", err
				}
				readFunc += fmt.Sprintln(findMultiCacheByFieldPlural.String())
				findMultiByFieldPlural, err := template.NewTemplate("findMultiByFieldPlural").Parse(FindMultiByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"dbName":           r.dbName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
					"whereField":       whereField,
				})
				if err != nil {
					return "", err
				}
				readFunc += fmt.Sprintln(findMultiByFieldPlural.String())
			}
		}
		// 唯一 && 字段数大于1
		if v.Unique && len(v.Columns) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var lowerFieldsJoin string
			var whereFields string
			for _, v := range v.Columns {
				upperFields += r.upperFieldName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.lowerFieldName(v), r.columnNameToDataType[v])
				lowerFieldsJoin += fmt.Sprintf("%s,", r.lowerFieldName(v))
				switch r.columnNameToDataType[v] {
				case "bool":
					whereFields += fmt.Sprintf("dao.%s.Is(%s),", r.upperFieldName(v), r.lowerFieldName(v))
				default:
					whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.upperFieldName(v), r.lowerFieldName(v))
				}
			}
			findOneCacheByFields, err := template.NewTemplate("findOneCacheByFields").Parse(FindOneCacheByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"dbName":            r.dbName,
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
			findOneByFields, err := template.NewTemplate("findOneByFields").Parse(FindOneByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"dbName":            r.dbName,
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
		if !v.Unique && len(v.Columns) == 1 {
			var whereField string
			columnNameToDataType := r.columnNameToDataType[v.Columns[0]]
			switch columnNameToDataType {
			case "bool":
				whereField += fmt.Sprintf("dao.%s.Is(%s),", r.upperFieldName(v.Columns[0]), r.lowerFieldName(v.Columns[0]))
			default:
				whereField += fmt.Sprintf("dao.%s.Eq(%s),", r.upperFieldName(v.Columns[0]), r.lowerFieldName(v.Columns[0]))
			}
			findMultiByField, err := template.NewTemplate("findMultiByField").Parse(FindMultiByField).Execute(map[string]any{
				"firstTableChar": r.firstTableChar,
				"dbName":         r.dbName,
				"upperTableName": r.upperTableName,
				"lowerTableName": r.lowerTableName,
				"upperField":     r.upperFieldName(v.Columns[0]),
				"lowerField":     r.lowerFieldName(v.Columns[0]),
				"dataType":       r.columnNameToDataType[v.Columns[0]],
				"whereField":     whereField,
			})
			if err != nil {
				return "", err
			}
			readFunc += fmt.Sprintln(findMultiByField.String())
			switch columnNameToDataType {
			case "bool":
			default:
				findMultiByFieldPlural, err := template.NewTemplate("findMultiByFieldPlural").Parse(FindMultiByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"dbName":           r.dbName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
					"whereField":       whereField,
				})
				if err != nil {
					return "", err
				}
				readFunc += fmt.Sprintln(findMultiByFieldPlural.String())
			}
		}
		// 不唯一 && 字段数大于1
		if !v.Unique && len(v.Columns) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var whereFields string
			for _, v := range v.Columns {
				upperFields += r.upperFieldName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.lowerFieldName(v), r.columnNameToDataType[v])
				switch r.columnNameToDataType[v] {
				case "bool":
					whereFields += fmt.Sprintf("dao.%s.Is(%s),", r.upperFieldName(v), r.lowerFieldName(v))
				default:
					whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.upperFieldName(v), r.lowerFieldName(v))
				}
			}
			findMultiByFields, err := template.NewTemplate("findMultiByFields").Parse(FindMultiByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"dbName":            r.dbName,
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
	findMultiByCustom, err := template.NewTemplate("FindMultiByCustom").Parse(FindMultiByCustom).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
	})
	if err != nil {
		return "", err
	}
	readFunc += fmt.Sprintln(findMultiByCustom.String())
	return readFunc, nil
}

// generateUpdateFunc
func (r *Repo) generateUpdateFunc() (string, error) {
	var updateFunc string
	var primaryKey string
	for _, v := range r.index {
		if v.PrimaryKey {
			primaryKey = v.Columns[0]
			break
		}
	}
	if primaryKey == "" {
		return "", nil
	}
	updateOneTpl, err := template.NewTemplate("UpdateOne").Parse(UpdateOne).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
		"upperField":     r.upperFieldName(primaryKey),
	})
	if err != nil {
		return "", err
	}
	updateFunc += fmt.Sprintln(updateOneTpl.String())
	updateOneByTx, err := template.NewTemplate("UpdateOneByTx").Parse(UpdateOneByTx).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
		"upperField":     r.upperFieldName(primaryKey),
	})
	if err != nil {
		return "", err
	}
	updateFunc += fmt.Sprintln(updateOneByTx.String())

	updateOneWithZero, err := template.NewTemplate("UpdateOneWithZero").Parse(UpdateOneWithZero).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
		"upperField":     r.upperFieldName(primaryKey),
	})
	if err != nil {
		return "", err
	}
	updateFunc += fmt.Sprintln(updateOneWithZero.String())
	updateOneWithZeroByTx, err := template.NewTemplate("UpdateOneWithZeroByTx").Parse(UpdateOneWithZeroByTx).Execute(map[string]any{
		"firstTableChar": r.firstTableChar,
		"dbName":         r.dbName,
		"upperTableName": r.upperTableName,
		"lowerTableName": r.lowerTableName,
		"upperField":     r.upperFieldName(primaryKey),
	})
	if err != nil {
		return "", err
	}
	updateFunc += fmt.Sprintln(updateOneWithZeroByTx.String())
	return updateFunc, nil
}

// generateDelFunc
func (r *Repo) generateDelFunc() (string, error) {
	var delMethods string
	var varCacheDelKeys string
	var haveUnique bool
	for _, v := range r.index {
		if r.checkDaoFieldType(v.Columns) {
			continue
		}
		if v.Unique {
			haveUnique = true
			var cacheField string
			cacheFieldsJoinSli := make([]string, 0)
			for _, column := range v.Columns {
				cacheField += r.upperFieldName(column)
				cacheFieldsJoinSli = append(cacheFieldsJoinSli, fmt.Sprintf("v.%s", r.upperFieldName(column)))
			}
			varCacheDelKeyTpl, err := template.NewTemplate("VarCacheDelKey").Parse(VarCacheDelKey).Execute(map[string]any{
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
		if v.Unique && len(v.Columns) == 1 {
			columnNameToDataType := r.columnNameToDataType[v.Columns[0]]
			switch columnNameToDataType {
			case "bool":
			default:
				deleteOneCacheByField, err := template.NewTemplate("DeleteOneCacheByField").Parse(DeleteOneCacheByField).Execute(map[string]any{
					"firstTableChar": r.firstTableChar,
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteOneCacheByField.String())
				deleteOneCacheByFieldTx, err := template.NewTemplate("DeleteOneCacheByFieldTx").Parse(DeleteOneCacheByFieldTx).Execute(map[string]any{
					"firstTableChar": r.firstTableChar,
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteOneCacheByFieldTx.String())
				deleteOneByField, err := template.NewTemplate("DeleteOneByField").Parse(DeleteOneByField).Execute(map[string]any{
					"firstTableChar": r.firstTableChar,
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteOneByField.String())
				deleteOneByFieldTx, err := template.NewTemplate("DeleteOneByFieldTx").Parse(DeleteOneByFieldTx).Execute(map[string]any{
					"firstTableChar": r.firstTableChar,
					"dbName":         r.dbName,
					"upperTableName": r.upperTableName,
					"lowerTableName": r.lowerTableName,
					"upperField":     r.upperFieldName(v.Columns[0]),
					"lowerField":     r.lowerFieldName(v.Columns[0]),
					"dataType":       r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteOneByFieldTx.String())
				deleteMultiCacheByFieldPlural, err := template.NewTemplate("DeleteMultiCacheByFieldPlural").Parse(DeleteMultiCacheByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteMultiCacheByFieldPlural.String())
				deleteMultiCacheByFieldPluralTx, err := template.NewTemplate("DeleteMultiCacheByFieldPluralTx").Parse(DeleteMultiCacheByFieldPluralTx).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteMultiCacheByFieldPluralTx.String())
				deleteMultiByFieldPlural, err := template.NewTemplate("DeleteMultiByFieldPlural").Parse(DeleteMultiByFieldPlural).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteMultiByFieldPlural.String())
				deleteMultiByFieldPluralTx, err := template.NewTemplate("DeleteMultiByFieldPluralTx").Parse(DeleteMultiByFieldPluralTx).Execute(map[string]any{
					"firstTableChar":   r.firstTableChar,
					"dbName":           r.dbName,
					"upperTableName":   r.upperTableName,
					"lowerTableName":   r.lowerTableName,
					"upperField":       r.upperFieldName(v.Columns[0]),
					"lowerField":       r.lowerFieldName(v.Columns[0]),
					"upperFieldPlural": r.plural(r.upperFieldName(v.Columns[0])),
					"lowerFieldPlural": r.plural(r.lowerFieldName(v.Columns[0])),
					"dataType":         r.columnNameToDataType[v.Columns[0]],
				})
				if err != nil {
					return "", err
				}
				delMethods += fmt.Sprintln(deleteMultiByFieldPluralTx.String())
			}
		}
		// 唯一 && 字段数大于1
		if v.Unique && len(v.Columns) > 1 {
			var upperFields string
			var fieldAndDataTypes string
			var whereFields string
			for _, v := range v.Columns {
				upperFields += r.upperFieldName(v)
				fieldAndDataTypes += fmt.Sprintf("%s %s,", r.lowerFieldName(v), r.columnNameToDataType[v])
				switch r.columnNameToDataType[v] {
				case "bool":
					whereFields += fmt.Sprintf("dao.%s.Is(%s),", r.upperFieldName(v), r.lowerFieldName(v))
				default:
					whereFields += fmt.Sprintf("dao.%s.Eq(%s),", r.upperFieldName(v), r.lowerFieldName(v))
				}
			}
			deleteOneCacheByFields, err := template.NewTemplate("DeleteOneCacheByFields").Parse(DeleteOneCacheByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.upperFieldName(v.Columns[0]),
				"lowerField":        r.lowerFieldName(v.Columns[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneCacheByFields.String())
			deleteOneCacheByFieldsTx, err := template.NewTemplate("DeleteOneCacheByFields").Parse(DeleteOneCacheByFieldsTx).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperField":        r.upperFieldName(v.Columns[0]),
				"lowerField":        r.lowerFieldName(v.Columns[0]),
				"upperFields":       upperFields,
				"dataType":          r.columnNameToDataType[v.Columns[0]],
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneCacheByFieldsTx.String())
			deleteOneByFields, err := template.NewTemplate("DeleteOneByFields").Parse(DeleteOneByFields).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"lowerField":        r.lowerFieldName(v.Columns[0]),
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneByFields.String())
			deleteOneByFieldsTx, err := template.NewTemplate("DeleteOneByFieldsTx").Parse(DeleteOneByFieldsTx).Execute(map[string]any{
				"firstTableChar":    r.firstTableChar,
				"dbName":            r.dbName,
				"upperTableName":    r.upperTableName,
				"lowerTableName":    r.lowerTableName,
				"upperFields":       upperFields,
				"lowerField":        r.lowerFieldName(v.Columns[0]),
				"fieldAndDataTypes": strings.Trim(fieldAndDataTypes, ","),
				"whereFields":       strings.Trim(whereFields, ","),
			})
			if err != nil {
				return "", err
			}
			delMethods += fmt.Sprintln(deleteOneByFieldsTx.String())
		}
		// 不唯一 && 字段数等于1
		if !v.Unique && len(v.Columns) == 1 {

		}
		// 不唯一 && 字段数大于1
		if !v.Unique && len(v.Columns) > 1 {

		}
	}
	if !haveUnique {
		return "", nil
	}
	deleteUniqueIndexCacheTpl, err := template.NewTemplate("DeleteUniqueIndexCache").Parse(DeleteUniqueIndexCache).Execute(map[string]any{
		"firstTableChar":  r.firstTableChar,
		"dbName":          r.dbName,
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

// upperFieldName 字段名称大写
func (r *Repo) upperFieldName(s string) string {
	return r.columnNameToName[s]
}

// lowerFieldName 字段名称小写
func (r *Repo) lowerFieldName(s string) string {
	str := r.upperFieldName(s)
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
	if token.Lookup(str).IsKeyword() || util.StrSliFind(KeyWords, str) {
		str = "_" + str
	}
	return str
}

// upperName 大写
func (r *Repo) upperName(s string) string {
	return r.gorm.NamingStrategy.SchemaName(s)
}

// lowerName 小写
func (r *Repo) lowerName(s string) string {
	str := r.upperName(s)
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

// plural 复数形式
func (r *Repo) plural(s string) string {
	str := inflection.Plural(s)
	if str == s {
		str += "plural"
	}
	return str
}

// checkDaoFieldType  检查字段状态
func (r *Repo) checkDaoFieldType(s []string) bool {
	for _, v := range s {
		if r.columnNameToFieldType[v] == "Field" {
			return true
		}
	}
	return false
}
