func New{{.upperTableName}}Repo(db *gorm.DB, redis *redis.Client) I{{.upperTableName}}Repo {
	return &{{.lowerTableName}}Repo{db: db, redis: redis}
}