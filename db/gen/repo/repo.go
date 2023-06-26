package repo

import (
	"fmt"
	"sync"

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
	fmt.Println(dbName)
	indexes, err := r.gorm.Migrator().GetIndexes(table)
	if err != nil {
		return err
	}
	fmt.Println(indexes)
	columnTypes, err := r.gorm.Migrator().ColumnTypes(table)
	if err != nil {
		return err
	}
	fmt.Println(columnTypes)
	pkgTpl, err := NewTemplate("pkg").Parse(Pkg).Execute(map[string]any{
		"lowerDbName": dbName,
	})
	if err != nil {
		return err
	}
	fmt.Println(pkgTpl)
	importTpl, err := NewTemplate("import").Parse(Import).Execute(map[string]any{
		"mod":          r.mod,
		"relativePath": r.relativePath,
		"lowerDbName":  dbName,
	})
	if err != nil {
		return err
	}
	fmt.Println(importTpl)
	return nil
}
