package gen

import (
	"fmt"

	"github.com/fzf-labs/fpkg/db/gen/repo"
	"gorm.io/gorm"
)

func GenerationRepo(gorm *gorm.DB) error {
	postgresqlModel := repo.NewPostgresqlModel(gorm)
	tables, err := postgresqlModel.FindAllTables()
	if err != nil {
		return err
	}
	fmt.Println(tables)
	return nil
}
