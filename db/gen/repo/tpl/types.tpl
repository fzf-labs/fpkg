type (
	I{{.upperTableName}}Repo interface{
		{{.methods}}
	}

	{{.upperTableName}}Repo struct {
		db    *gorm.DB
		rockscache *rockscache.Client
	}
)
