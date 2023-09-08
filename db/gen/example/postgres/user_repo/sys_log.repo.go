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

var _ ISysLogRepo = (*SysLogRepo)(nil)

var (
	cacheSysLogByIDPrefix = "DBCache:user:SysLogByID"
)

type (
	ISysLogRepo interface {
		// CreateOne 创建一条数据
		CreateOne(ctx context.Context, data *user_model.SysLog) error
		// CreateBatch 批量创建数据
		CreateBatch(ctx context.Context, data []*user_model.SysLog, batchSize int) error
		// UpdateOne 更新一条数据
		UpdateOne(ctx context.Context, data *user_model.SysLog) error
		// FindOneCacheByID 根据ID查询一条数据并设置缓存
		FindOneCacheByID(ctx context.Context, ID string) (*user_model.SysLog, error)
		// FindOneByID 根据ID查询一条数据
		FindOneByID(ctx context.Context, ID string) (*user_model.SysLog, error)
		// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
		FindMultiCacheByIDS(ctx context.Context, IDS []string) ([]*user_model.SysLog, error)
		// FindMultiByIDS 根据IDS查询多条数据
		FindMultiByIDS(ctx context.Context, IDS []string) ([]*user_model.SysLog, error)
		// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
		DeleteOneCacheByID(ctx context.Context, ID string) error
		// DeleteOneByID 根据ID删除一条数据
		DeleteOneByID(ctx context.Context, ID string) error
		// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
		DeleteMultiCacheByIDS(ctx context.Context, IDS []string) error
		// DeleteMultiByIDS 根据IDS删除多条数据
		DeleteMultiByIDS(ctx context.Context, IDS []string) error
		// DeleteUniqueIndexCache 删除唯一索引存在的缓存
		DeleteUniqueIndexCache(ctx context.Context, data []*user_model.SysLog) error
	}
	ISysLogCache interface {
		Key(fields ...any) string
		Fetch(ctx context.Context, key string, KvFn func() (string, error)) (string, error)
		FetchBatch(ctx context.Context, keys []string, KvFn func(miss []string) (map[string]string, error)) (map[string]string, error)
		Del(ctx context.Context, key string) error
		DelBatch(ctx context.Context, keys []string) error
	}
	SysLogRepo struct {
		db    *gorm.DB
		cache ISysLogCache
	}
)

func NewSysLogRepo(db *gorm.DB, cache ISysLogCache) *SysLogRepo {
	return &SysLogRepo{
		db:    db,
		cache: cache,
	}
}

// CreateOne 创建一条数据
func (r *SysLogRepo) CreateOne(ctx context.Context, data *user_model.SysLog) error {
	dao := user_dao.Use(r.db).SysLog
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}

// CreateBatch 批量创建数据
func (r *SysLogRepo) CreateBatch(ctx context.Context, data []*user_model.SysLog, batchSize int) error {
	dao := user_dao.Use(r.db).SysLog
	err := dao.WithContext(ctx).CreateInBatches(data, batchSize)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOne 更新一条数据
func (r *SysLogRepo) UpdateOne(ctx context.Context, data *user_model.SysLog) error {
	dao := user_dao.Use(r.db).SysLog
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	err = r.DeleteUniqueIndexCache(ctx, []*user_model.SysLog{data})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
func (r *SysLogRepo) DeleteOneCacheByID(ctx context.Context, ID string) error {
	dao := user_dao.Use(r.db).SysLog
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
	err = r.DeleteUniqueIndexCache(ctx, []*user_model.SysLog{first})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByID 根据ID删除一条数据
func (r *SysLogRepo) DeleteOneByID(ctx context.Context, ID string) error {
	dao := user_dao.Use(r.db).SysLog
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
func (r *SysLogRepo) DeleteMultiCacheByIDS(ctx context.Context, IDS []string) error {
	dao := user_dao.Use(r.db).SysLog
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
func (r *SysLogRepo) DeleteMultiByIDS(ctx context.Context, IDS []string) error {
	dao := user_dao.Use(r.db).SysLog
	_, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteUniqueIndexCache 删除唯一索引存在的缓存
func (r *SysLogRepo) DeleteUniqueIndexCache(ctx context.Context, data []*user_model.SysLog) error {
	keys := make([]string, 0)
	for _, v := range data {
		keys = append(keys, r.cache.Key(ctx, cacheSysLogByIDPrefix, v.ID))

	}
	err := r.cache.DelBatch(ctx, keys)
	if err != nil {
		return err
	}
	return nil
}

// FindOneCacheByID 根据ID查询一条数据并设置缓存
func (r *SysLogRepo) FindOneCacheByID(ctx context.Context, ID string) (*user_model.SysLog, error) {
	resp := new(user_model.SysLog)
	key := r.cache.Key(ctx, cacheSysLogByIDPrefix, ID)
	cacheValue, err := r.cache.Fetch(ctx, key, func() (string, error) {
		dao := user_dao.Use(r.db).SysLog
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
func (r *SysLogRepo) FindOneByID(ctx context.Context, ID string) (*user_model.SysLog, error) {
	dao := user_dao.Use(r.db).SysLog
	result, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return result, nil
}

// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
func (r *SysLogRepo) FindMultiCacheByIDS(ctx context.Context, IDS []string) ([]*user_model.SysLog, error) {
	resp := make([]*user_model.SysLog, 0)
	keys := make([]string, 0)
	keyToParam := make(map[string]string)
	for _, v := range IDS {
		key := r.cache.Key(ctx, cacheSysLogByIDPrefix, v)
		keys = append(keys, key)
		keyToParam[key] = v
	}
	cacheValue, err := r.cache.FetchBatch(ctx, keys, func(miss []string) (map[string]string, error) {
		params := make([]string, 0)
		for _, v := range miss {
			params = append(params, keyToParam[v])
		}
		dao := user_dao.Use(r.db).SysLog
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
			value[r.cache.Key(ctx, cacheSysLogByIDPrefix, v.ID)] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new(user_model.SysLog)
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}

// FindMultiByIDS 根据IDS查询多条数据
func (r *SysLogRepo) FindMultiByIDS(ctx context.Context, IDS []string) ([]*user_model.SysLog, error) {
	dao := user_dao.Use(r.db).SysLog
	result, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}
