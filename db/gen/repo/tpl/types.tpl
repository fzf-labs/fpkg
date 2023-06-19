type (
	I{{.upperTableName}}Repo interface{
		{{.methods}}
	}

	{{.lowerTableName}}Repo struct {
		db    *gorm.DB
		redis *redis.Client
	}
)
