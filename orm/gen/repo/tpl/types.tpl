type (
	I{{.upperTableName}}Repo interface{
		{{.methods}}
	}
	{{.upperTableName}}Repo struct {
		db    *gorm.DB
		cache cache.IDBCache
	}
)
