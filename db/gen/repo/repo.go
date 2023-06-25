package repo

import (
	"sync"

	"gorm.io/gorm"
)

func GenerationRepo(gorm *gorm.DB) error {
	postgresqlModel := NewPostgresqlModel(gorm)
	tables, err := postgresqlModel.FindAllTables()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, table := range tables {
		wg.Add(1)
		go func(table string) {
			_ = GenerationTable(gorm, table)
			wg.Done()
		}(table)
	}
	wg.Wait()
	return nil
}

func GenerationTable(gorm *gorm.DB, table string) error {
	//postgresqlModel := NewPostgresqlModel(gorm)
	//dbName = gorm.Migrator().CurrentDatabase()
	//dbColumns, err := postgresqlModel.FindColumns(table)
	//if err != nil {
	//	return err
	//}
	//dbIndices, err := postgresqlModel.FindIndex(table)
	//if err != nil {
	//	return err
	//}
	return nil
}
