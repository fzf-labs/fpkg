func New{{.upperTableName}}Repo(db *gorm.DB,cache cache.IDBCache) *{{.upperTableName}}Repo {
	return &{{.upperTableName}}Repo{
		db: db,
		cache:  cache,
	}
}