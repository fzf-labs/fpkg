func (u *userRepo) FindOneById(ctx context.Context, id string) (*user_model.User, error) {
	resp := new(user_model.User)
	cache := CacheById.NewSingleKey(u.redis)
	cacheValue, err := cache.SingleCache(ctx, id, func() (string, error) {
		userDao := user_dao.Use(u.db).User
		result, err := userDao.WithContext(ctx).Where(userDao.ID.Eq(id)).First()
		if err != nil && err != gorm.ErrRecordNotFound {
			return "", err
		}
		marshal, err := json.Marshal(result)
		if err != nil {
			return "", err
		}
		return string(marshal), nil
	})
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(cacheValue), resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}