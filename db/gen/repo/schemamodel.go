package repo

type (
	DbTable struct {
		DBName    string
		TableName string
	}

	DbColumn struct {
		Num           int    `gorm:"is_nullable"`
		ColumnName    string `gorm:"is_nullable"`
		DataType      string `gorm:"is_nullable"`
		IsNullAble    bool   `gorm:"is_nullable"`
		ColumnDefault string `gorm:"column_default"`
		Comment       string `gorm:"comment"`
	}

	DbIndex struct {
		IndexName      string
		IndexAlgorithm string
		IsUnique       bool
		ColumnName     string
	}
)
