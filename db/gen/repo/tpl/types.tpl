type (
	I{{.upperTableName}}Repo interface{
		{{.methods}}
	}
	I{{.upperTableName}}Cache interface {
		Key(ctx context.Context, keys ...any) string
		Take(ctx context.Context, keys []string, KvFn func() (map[string]string, error)) (map[string]string, error)
		Del(ctx context.Context, keys []string) error
	}
	{{.upperTableName}}Repo struct {
		db    *gorm.DB
		cache I{{.upperTableName}}Cache
	}
)
