// Code generated by gen/repo. DO NOT EDIT.
// Code generated by gen/repo. DO NOT EDIT.
// Code generated by gen/repo. DO NOT EDIT.

package user_repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fzf-labs/fpkg/db/gen/example/postgres/user_dao"
	"github.com/fzf-labs/fpkg/db/gen/example/postgres/user_model"
	"gorm.io/gorm"
)

var _ ISysDictRepo = (*SysDictRepo)(nil)

var (
	cacheSysDictByIDPrefix = "DBCache:user:SysDictByID"
)

type (
	ISysDictRepo interface {
		// CreateOne 创建一条数据
		CreateOne(ctx context.Context, data *user_model.SysDict) error
		// CreateBatch 批量创建数据
		CreateBatch(ctx context.Context, data []*user_model.SysDict, batchSize int) error
		// UpdateOne 更新一条数据
		UpdateOne(ctx context.Context, data *user_model.SysDict) error
		// FindOneCacheByID 根据ID查询一条数据并设置缓存
		FindOneCacheByID(ctx context.Context, ID string) (*user_model.SysDict, error)
		// FindOneByID 根据ID查询一条数据
		FindOneByID(ctx context.Context, ID string) (*user_model.SysDict, error)
		// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
		FindMultiCacheByIDS(ctx context.Context, IDS []string) ([]*user_model.SysDict, error)
		// FindMultiByIDS 根据IDS查询多条数据
		FindMultiByIDS(ctx context.Context, IDS []string) ([]*user_model.SysDict, error)
		// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
		DeleteOneCacheByID(ctx context.Context, ID string) error
		// DeleteOneByID 根据ID删除一条数据
		DeleteOneByID(ctx context.Context, ID string) error
		// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
		DeleteMultiCacheByIDS(ctx context.Context, IDS []string) error
		// DeleteMultiByIDS 根据IDS删除多条数据
		DeleteMultiByIDS(ctx context.Context, IDS []string) error
		// DeleteUniqueIndexCache 删除唯一索引存在的缓存
		DeleteUniqueIndexCache(ctx context.Context, data []*user_model.SysDict) error
	}
	ISysDictCache interface {
		Key(fields ...any) string
		Fetch(ctx context.Context, key string, KvFn func() (string, error)) (string, error)
		FetchBatch(ctx context.Context, keys []string, KvFn func(miss []string) (map[string]string, error)) (map[string]string, error)
		Del(ctx context.Context, key string) error
		DelBatch(ctx context.Context, keys []string) error
	}
	SysDictRepo struct {
		db    *gorm.DB
		cache ISysDictCache
	}
)

func NewSysDictRepo(db *gorm.DB, cache ISysDictCache) *SysDictRepo {
	return &SysDictRepo{
		db:    db,
		cache: cache,
	}
}

// CreateOne 创建一条数据
func (r *SysDictRepo) CreateOne(ctx context.Context, data *user_model.SysDict) error {
	dao := user_dao.Use(r.db).SysDict
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}

// CreateBatch 批量创建数据
func (r *SysDictRepo) CreateBatch(ctx context.Context, data []*user_model.SysDict, batchSize int) error {
	dao := user_dao.Use(r.db).SysDict
	err := dao.WithContext(ctx).CreateInBatches(data, batchSize)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOne 更新一条数据
func (r *SysDictRepo) UpdateOne(ctx context.Context, data *user_model.SysDict) error {
	dao := user_dao.Use(r.db).SysDict
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	err = r.DeleteUniqueIndexCache(ctx, []*user_model.SysDict{data})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
func (r *SysDictRepo) DeleteOneCacheByID(ctx context.Context, ID string) error {
	dao := user_dao.Use(r.db).SysDict
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
	err = r.DeleteUniqueIndexCache(ctx, []*user_model.SysDict{first})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByID 根据ID删除一条数据
func (r *SysDictRepo) DeleteOneByID(ctx context.Context, ID string) error {
	dao := user_dao.Use(r.db).SysDict
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
func (r *SysDictRepo) DeleteMultiCacheByIDS(ctx context.Context, IDS []string) error {
	dao := user_dao.Use(r.db).SysDict
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

// DeleteMultiByIDS 根据IDS删除多条数据
func (r *SysDictRepo) DeleteMultiByIDS(ctx context.Context, IDS []string) error {
	dao := user_dao.Use(r.db).SysDict
	_, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteUniqueIndexCache 删除唯一索引存在的缓存
func (r *SysDictRepo) DeleteUniqueIndexCache(ctx context.Context, data []*user_model.SysDict) error {
	keys := make([]string, 0)
	for _, v := range data {
		keys = append(keys, r.cache.Key(ctx, cacheSysDictByIDPrefix, v.ID))

	}
	err := r.cache.DelBatch(ctx, keys)
	if err != nil {
		return err
	}
	return nil
}

// FindOneCacheByID 根据ID查询一条数据并设置缓存
func (r *SysDictRepo) FindOneCacheByID(ctx context.Context, ID string) (*user_model.SysDict, error) {
	resp := new(user_model.SysDict)
	key := r.cache.Key(ctx, cacheSysDictByIDPrefix, ID)
	cacheValue, err := r.cache.Fetch(ctx, key, func() (string, error) {
		dao := user_dao.Use(r.db).SysDict
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

// FindOneByID 根据ID查询一条数据
func (r *SysDictRepo) FindOneByID(ctx context.Context, ID string) (*user_model.SysDict, error) {
	dao := user_dao.Use(r.db).SysDict
	result, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return result, nil
}

// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
func (r *SysDictRepo) FindMultiCacheByIDS(ctx context.Context, IDS []string) ([]*user_model.SysDict, error) {
	resp := make([]*user_model.SysDict, 0)
	keys := make([]string, 0)
	keyToParam := make(map[string]string)
	for _, v := range IDS {
		key := r.cache.Key(ctx, cacheSysDictByIDPrefix, v)
		keys = append(keys, key)
		keyToParam[key] = v
	}
	cacheValue, err := r.cache.FetchBatch(ctx, keys, func(miss []string) (map[string]string, error) {
		params := make([]string, 0)
		for _, v := range miss {
			params = append(params, keyToParam[v])
		}
		dao := user_dao.Use(r.db).SysDict
		result, err := dao.WithContext(ctx).Where(dao.ID.In(params...)).Find()
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		value := make(map[string]string)
		for _, v := range miss {
			value[v] = ""
		}
		for _, v := range result {
			marshal, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			value[r.cache.Key(ctx, cacheSysDictByIDPrefix, v.ID)] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new(user_model.SysDict)
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}

// FindMultiByIDS 根据IDS查询多条数据
func (r *SysDictRepo) FindMultiByIDS(ctx context.Context, IDS []string) ([]*user_model.SysDict, error) {
	dao := user_dao.Use(r.db).SysDict
	result, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}
