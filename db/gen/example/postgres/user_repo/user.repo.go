// Code generated by gen/repo. DO NOT EDIT.
// Code generated by gen/repo. DO NOT EDIT.
// Code generated by gen/repo. DO NOT EDIT.

package user_repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fzf-labs/fpkg/cache/cachekey"
	"github.com/fzf-labs/fpkg/conv"
	"github.com/fzf-labs/fpkg/db/gen/example/postgres/user_dao"
	"github.com/fzf-labs/fpkg/db/gen/example/postgres/user_model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ IUserRepo = (*UserRepo)(nil)

var (
	// 缓存管理器
	cacheKeyUserManage = cachekey.NewKeyManage("UserRepo")
	// 只针对唯一索引做缓存
	CacheUserByID = cacheKeyUserManage.AddKey("CacheUserByID", time.Hour*24, "CacheUserByID")
)

type (
	IUserRepo interface {
		// CreateOne 创建一条数据
		CreateOne(ctx context.Context, data *user_model.User) error
		// UpdateOne 更新一条数据
		UpdateOne(ctx context.Context, data *user_model.User) error
		// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
		DeleteOneCacheByID(ctx context.Context, ID int64) error
		// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
		DeleteMultiCacheByIDS(ctx context.Context, IDS []int64) error
		// DeleteUniqueIndexCache 删除唯一索引存在的缓存
		DeleteUniqueIndexCache(ctx context.Context, data []*user_model.User) error
		// FindOneCacheByID 根据ID查询一条数据并设置缓存
		FindOneCacheByID(ctx context.Context, ID int64) (*user_model.User, error)
		// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
		FindMultiCacheByIDS(ctx context.Context, IDS []int64) ([]*user_model.User, error)
	}

	UserRepo struct {
		db    *gorm.DB
		redis *redis.Client
	}
)

func NewUserRepo(db *gorm.DB, redis *redis.Client) *UserRepo {
	return &UserRepo{db: db, redis: redis}
}

// CreateOne 创建一条数据
func (r *UserRepo) CreateOne(ctx context.Context, data *user_model.User) error {
	dao := user_dao.Use(r.db).User
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOne 更新一条数据
func (r *UserRepo) UpdateOne(ctx context.Context, data *user_model.User) error {
	dao := user_dao.Use(r.db).User
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	err = r.DeleteUniqueIndexCache(ctx, []*user_model.User{data})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
func (r *UserRepo) DeleteOneCacheByID(ctx context.Context, ID int64) error {
	dao := user_dao.Use(r.db).User
	first, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if first == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	err = r.DeleteUniqueIndexCache(ctx, []*user_model.User{first})
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
func (r *UserRepo) DeleteMultiCacheByIDS(ctx context.Context, IDS []int64) error {
	dao := user_dao.Use(r.db).User
	list, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Find()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	err = r.DeleteUniqueIndexCache(ctx, list)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUniqueIndexCache 删除唯一索引存在的缓存
func (r *UserRepo) DeleteUniqueIndexCache(ctx context.Context, data []*user_model.User) error {
	var err error
	cacheUserByID := CacheUserByID.NewSingleKey(r.redis)

	for _, v := range data {
		err = cacheUserByID.SingleCacheDel(ctx, cacheUserByID.BuildKey(v.ID))
		if err != nil {
			return err
		}

	}
	return nil
}

// FindOneCacheByID 根据ID查询一条数据并设置缓存
func (r *UserRepo) FindOneCacheByID(ctx context.Context, ID int64) (*user_model.User, error) {
	resp := new(user_model.User)
	cache := CacheUserByID.NewSingleKey(r.redis)
	cacheValue, err := cache.SingleCache(ctx, conv.String(ID), func() (string, error) {
		dao := user_dao.Use(r.db).User
		result, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).First()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
func (r *UserRepo) FindMultiCacheByIDS(ctx context.Context, IDS []int64) ([]*user_model.User, error) {
	resp := make([]*user_model.User, 0)
	cacheKey := CacheUserByID.NewBatchKey(r.redis)
	batchKeys := make([]string, 0)
	for _, v := range IDS {
		batchKeys = append(batchKeys, conv.String(v))
	}
	cacheValue, err := cacheKey.BatchKeyCache(ctx, batchKeys, func() (map[string]string, error) {
		dao := user_dao.Use(r.db).User
		result, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Find()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		value := make(map[string]string)
		for _, v := range result {
			marshal, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			value[conv.String(v.ID)] = string(marshal)
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
