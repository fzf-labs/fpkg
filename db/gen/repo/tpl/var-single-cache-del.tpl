		err = cache{{.upperTableName}}By{{.cacheField}}.SingleCacheDel(ctx, cache{{.upperTableName}}By{{.cacheField}}.BuildKey({{.cacheFieldsJoin}}))
		if err != nil {
			return err
		}