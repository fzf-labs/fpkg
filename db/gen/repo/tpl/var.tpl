var _ I{{.upperTableName}}Repo = (*{{.lowerTableName}}Repo)(nil)

var (
	// 缓存管理器
	cacheKey{{.upperTableName}}Manage = cachekey.NewKeyManage("{{.lowerTableName}}Repo")
	// 只针对唯一索引做缓存
    {{.cacheKeys}}
)