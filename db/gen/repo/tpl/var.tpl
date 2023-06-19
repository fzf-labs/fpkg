var _ I{{.upperTableName}}Repo = (*{{.lowerTableName}}Repo)(nil)

var (
	cacheKeyManage = cachekey.NewKeyManage("{{.lowerTableName}}_repo")
    {{.cacheKeys}}
)