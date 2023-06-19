package user_repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fzf-labs/fpkg/cache/cachekey"
	"github.com/fzf-labs/fpkg/db/gen/example/postgres/user_dao"
	"github.com/fzf-labs/fpkg/db/gen/example/postgres/user_model"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var _ IUserRepo = (*userRepo)(nil)

var (
	cacheKeyManage = cachekey.NewKeyManage("user_repo")
	CacheById      = cacheKeyManage.AddKey("CacheById", time.Hour, "CacheById")
	CacheByPhone   = cacheKeyManage.AddKey("CacheByPhone", time.Hour, "CacheById")
	CacheByEmail   = cacheKeyManage.AddKey("CacheByEmail", time.Hour, "CacheById")
)

type (
	IUserRepo interface {
		CreateOne(ctx context.Context, data *user_model.User) error
		UpdateOne(ctx context.Context, data *user_model.User) error

		FindOneById(ctx context.Context, id string) (*user_model.User, error)
		DeleteOneById(ctx context.Context, id string) error

		FindMultiByIds(ctx context.Context, ids []string) ([]*user_model.User, error)
		DeleteMultiByIds(ctx context.Context, ids []string) error

		FindOneByPhone(ctx context.Context, phone string) (*user_model.User, error)
		DeleteOneByPhone(ctx context.Context, phone string) error

		FindMultiByPhones(ctx context.Context, phones []string) ([]*user_model.User, error)
		DeleteMultiByPhones(ctx context.Context, phones []string) error
	}
	userRepo struct {
		db    *gorm.DB
		redis *redis.Client
	}
)

func NewUserRepo(db *gorm.DB, redis *redis.Client) IUserRepo {
	return &userRepo{db: db, redis: redis}
}

func (u *userRepo) CreateOne(ctx context.Context, data *user_model.User) error {
	userDao := user_dao.Use(u.db).User
	err := userDao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) UpdateOne(ctx context.Context, data *user_model.User) error {
	userDao := user_dao.Use(u.db).User
	_, err := userDao.WithContext(ctx).Where(userDao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	cache := CacheById.NewSingleKey(u.redis)
	err = cache.SingleCacheDel(ctx, data.ID)
	if err != nil {
		return err
	}
	return nil
}

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

func (u *userRepo) FindMultiByIds(ctx context.Context, ids []string) ([]*user_model.User, error) {
	resp := make([]*user_model.User, 0)
	cacheKey := CacheById.NewBatchKey(u.redis)
	cacheValue, err := cacheKey.BatchKeyCache(ctx, ids, func() (map[string]string, error) {
		userDao := user_dao.Use(u.db).User
		result, err := userDao.WithContext(ctx).Where(userDao.ID.In(ids...)).Find()
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		value := make(map[string]string)
		for _, v := range result {
			marshal, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			value[v.ID] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new(user_model.User)
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}
func (u *userRepo) DeleteOneById(ctx context.Context, id string) error {
	userDao := user_dao.Use(u.db).User
	_, err := userDao.WithContext(ctx).Where(userDao.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}
	cache := CacheById.NewSingleKey(u.redis)
	err = cache.SingleCacheDel(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (u *userRepo) DeleteMultiByIds(ctx context.Context, ids []string) error {
	userDao := user_dao.Use(u.db).User
	_, err := userDao.WithContext(ctx).Where(userDao.ID.In(ids...)).Delete()
	if err != nil {
		return err
	}
	cacheKey := CacheById.NewBatchKey(u.redis)
	err = cacheKey.BatchKeyCacheDel(ctx, ids)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) FindOneByPhone(ctx context.Context, phone string) (*user_model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) FindMultiByPhones(ctx context.Context, phones []string) ([]*user_model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) DeleteOneByPhone(ctx context.Context, phone string) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) DeleteMultiByPhones(ctx context.Context, phones []string) error {
	//TODO implement me
	panic("implement me")
}
