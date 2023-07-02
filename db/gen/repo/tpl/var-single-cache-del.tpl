		err = cacheBy{{.cacheField}}.SingleCacheDel(ctx, cacheBy{{.cacheField}}.BuildKey({{.cacheFieldsJoin}}))
		if err != nil {
			return err
		}