package repo

import (
	"fmt"
	"sync"

	"github.com/fzf-labs/fpkg/util/jsonutil"
	"gorm.io/gorm"
)

type Repo struct {
	gorm         *gorm.DB
	mod          string
	relativePath string
}

func NewRepo(gorm *gorm.DB, mod string, relativePath string) *Repo {
	return &Repo{gorm: gorm, mod: mod, relativePath: relativePath}
}

func (r *Repo) GenerationRepo() error {
	//获取表
	tables, err := r.gorm.Migrator().GetTables()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, table := range tables {
		wg.Add(1)
		go func(table string) {
			_ = r.GenerationTable(table)
			wg.Done()
		}(table)
	}
	wg.Wait()
	return nil
}

func (r *Repo) GenerationTable(table string) error {
	dbName := r.gorm.Migrator().CurrentDatabase()
	jsonutil.Dump(dbName)
	indexes, err := r.gorm.Migrator().GetIndexes(table)
	if err != nil {
		return err
	}
	jsonutil.Dump(indexes)
	columnTypes, err := r.gorm.Migrator().ColumnTypes(table)
	if err != nil {
		return err
	}
	jsonutil.Dump(columnTypes)
	pkgTpl, err := NewTemplate("Pkg").Parse(Pkg).Execute(map[string]any{
		"lowerDbName": dbName,
	})
	if err != nil {
		return err
	}
	fmt.Println(pkgTpl)
	importTpl, err := NewTemplate("Import").Parse(Import).Execute(map[string]any{
		"mod":          r.mod,
		"relativePath": r.relativePath,
		"lowerDbName":  dbName,
	})
	if err != nil {
		return err
	}
	fmt.Println(importTpl)
	upperTableName := UpperName(table)
	lowerTableName := LowerName(table)
	var cacheKeys string
	var methods string
	interfaceCreateOneTpl, err := NewTemplate("InterfaceCreateOne").Parse(InterfaceCreateOne).Execute(map[string]any{
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
	})
	if err != nil {
		return err
	}
	interfaceUpdateOneTpl, err := NewTemplate("InterfaceUpdateOne").Parse(InterfaceUpdateOne).Execute(map[string]any{
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
	})

	methods += interfaceCreateOneTpl.String()
	methods += interfaceUpdateOneTpl.String()
	for _, index := range indexes {
		//唯一性索引
		unique, _ := index.Unique()
		if unique {
			var cacheField string
			for _, column := range index.Columns() {
				cacheField += UpperName(column)
			}
			varCacheTpl, err := NewTemplate("VarCache").Parse(VarCache).Execute(map[string]any{
				"cacheField": cacheField,
			})
			if err != nil {
				return err
			}
			cacheKeys += varCacheTpl.String()

			if len(index.Columns()) > 1 {
				interfaceDeleteMultiByFieldComplex, err := NewTemplate("InterfaceDeleteMultiByFieldComplex").Parse(InterfaceDeleteMultiByFieldComplex).Execute(map[string]any{
					"upperTableName": upperTableName,
					"lowerTableName": lowerTableName,
				})
				if err != nil {
					return err
				}
				fmt.Println(interfaceDeleteMultiByFieldComplex)
			} else {
				interfaceFindOneCacheByField, err := NewTemplate("InterfaceFindOneCacheByField").Parse(InterfaceFindOneCacheByField).Execute(map[string]any{
					"upperTableName": upperTableName,
					"lowerTableName": lowerTableName,
					"upperField":     UpperName(index.Columns()[0]),
					"lowerField":     LowerName(index.Columns()[0]),
					"dataType":       lowerTableName,
				})
				fmt.Println(interfaceFindOneCacheByField)

				if err != nil {
					return err
				}
				interfaceFindMultiCacheByFieldComplex, err := NewTemplate("InterfaceFindMultiCacheByFieldComplex").Parse(InterfaceFindMultiCacheByFieldComplex).Execute(map[string]any{
					"upperTableName": upperTableName,
					"lowerTableName": lowerTableName,
				})
				if err != nil {
					return err
				}
				fmt.Println(interfaceFindMultiCacheByFieldComplex)
			}

		} else {

		}
	}
	varTpl, err := NewTemplate("Var").Parse(Var).Execute(map[string]any{
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
		"cacheKeys":      cacheKeys,
	})
	if err != nil {
		return err
	}
	fmt.Println(varTpl)
	typesTpl, err := NewTemplate("Types").Parse(Types).Execute(map[string]any{
		"upperTableName": upperTableName,
		"lowerTableName": lowerTableName,
		"methods":        methods,
	})
	if err != nil {
		return err
	}
	fmt.Println(typesTpl)
	return nil
}
