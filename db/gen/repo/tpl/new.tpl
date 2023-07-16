func New{{.upperTableName}}Repo(db *gorm.DB, redis *redis.Client) I{{.upperTableName}}Repo {
	return &{{.upperTableName}}Repo{db: db, redis: redis}
}