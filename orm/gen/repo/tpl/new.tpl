func New{{.upperTableName}}Repo(db *gorm.DB,cache I{{.upperTableName}}Cache) *{{.upperTableName}}Repo {
	return &{{.upperTableName}}Repo{
		db: db,
		cache:  cache,
	}
}