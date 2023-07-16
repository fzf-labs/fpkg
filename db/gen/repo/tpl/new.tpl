func New{{.upperTableName}}Repo(db *gorm.DB, redis *redis.Client) *{{.upperTableName}}Repo {
	return &{{.upperTableName}}Repo{db: db, redis: redis}
}