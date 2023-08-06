func New{{.upperTableName}}Repo(db *gorm.DB,rockscache *rockscache.Client) *{{.upperTableName}}Repo {
	return &{{.upperTableName}}Repo{
		db:         db,
		rockscache: rockscache,
	}
}