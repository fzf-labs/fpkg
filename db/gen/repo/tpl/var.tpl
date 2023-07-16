var _ I{{.upperTableName}}Repo = (*{{.upperTableName}}Repo)(nil)

var (
	// 缓存管理器
	cacheKey{{.upperTableName}}Manage = cachekey.NewKeyManage("{{.upperTableName}}Repo")
	// 只针对唯一索引做缓存
    {{.cacheKeys}}
)