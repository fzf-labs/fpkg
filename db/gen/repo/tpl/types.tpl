type (
	I{{.upperTableName}}Repo interface{
		{{.methods}}
	}
	I{{.upperTableName}}Cache interface {
        Key(fields ...any) string
        Fetch(ctx context.Context, key string, fn func() (string, error)) (string, error)
        FetchBatch(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error)) (map[string]string, error)
        Del(ctx context.Context, key string) error
        DelBatch(ctx context.Context, keys []string) error
	}
	{{.upperTableName}}Repo struct {
		db    *gorm.DB
		cache I{{.upperTableName}}Cache
	}
)
