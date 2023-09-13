// Code generated by gen/repo. DO NOT EDIT.
// Code generated by gen/repo. DO NOT EDIT.
// Code generated by gen/repo. DO NOT EDIT.

package gorm_gen_repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/fzf-labs/fpkg/orm"
	"github.com/fzf-labs/fpkg/orm/gen/example/postgres/gorm_gen_dao"
	"github.com/fzf-labs/fpkg/orm/gen/example/postgres/gorm_gen_model"
	"gorm.io/gorm"
)

var _ ISysAdminRepo = (*SysAdminRepo)(nil)

var (
	cacheSysAdminByUsernamePrefix = "DBCache:gorm_gen:SysAdminByUsername"
	cacheSysAdminByIDPrefix       = "DBCache:gorm_gen:SysAdminByID"
)

type (
	ISysAdminRepo interface {
		// CreateOne 创建一条数据
		CreateOne(ctx context.Context, data *gorm_gen_model.SysAdmin) error
		// CreateOneByTx 创建一条数据(事务)
		CreateOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.SysAdmin) error
		// CreateBatch 批量创建数据
		CreateBatch(ctx context.Context, data []*gorm_gen_model.SysAdmin, batchSize int) error
		// UpdateOne 更新一条数据
		UpdateOne(ctx context.Context, data *gorm_gen_model.SysAdmin) error
		// UpdateOne 更新一条数据(事务)
		UpdateOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.SysAdmin) error
		// FindOneCacheByUsername 根据username查询一条数据并设置缓存
		FindOneCacheByUsername(ctx context.Context, username string) (*gorm_gen_model.SysAdmin, error)
		// FindOneByUsername 根据username查询一条数据
		FindOneByUsername(ctx context.Context, username string) (*gorm_gen_model.SysAdmin, error)
		// FindMultiCacheByUsernames 根据usernames查询多条数据并设置缓存
		FindMultiCacheByUsernames(ctx context.Context, usernames []string) ([]*gorm_gen_model.SysAdmin, error)
		// FindMultiByUsernames 根据usernames查询多条数据
		FindMultiByUsernames(ctx context.Context, usernames []string) ([]*gorm_gen_model.SysAdmin, error)
		// FindOneCacheByID 根据ID查询一条数据并设置缓存
		FindOneCacheByID(ctx context.Context, ID string) (*gorm_gen_model.SysAdmin, error)
		// FindOneByID 根据ID查询一条数据
		FindOneByID(ctx context.Context, ID string) (*gorm_gen_model.SysAdmin, error)
		// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
		FindMultiCacheByIDS(ctx context.Context, IDS []string) ([]*gorm_gen_model.SysAdmin, error)
		// FindMultiByIDS 根据IDS查询多条数据
		FindMultiByIDS(ctx context.Context, IDS []string) ([]*gorm_gen_model.SysAdmin, error)
		// FindMultiByPaginator 查询分页数据(通用)
		FindMultiByPaginator(ctx context.Context, params *orm.PaginatorParams) ([]*gorm_gen_model.SysAdmin, int64, error)
		// DeleteOneCacheByUsername 根据username删除一条数据并清理缓存
		DeleteOneCacheByUsername(ctx context.Context, username string) error
		// DeleteOneCacheByUsername 根据username删除一条数据并清理缓存
		DeleteOneCacheByUsernameTx(ctx context.Context, tx *gorm_gen_dao.Query, username string) error
		// DeleteOneByUsername 根据username删除一条数据
		DeleteOneByUsername(ctx context.Context, username string) error
		// DeleteOneByUsername 根据username删除一条数据
		DeleteOneByUsernameTx(ctx context.Context, tx *gorm_gen_dao.Query, username string) error
		// DeleteMultiCacheByUsernames 根据Usernames删除多条数据并清理缓存
		DeleteMultiCacheByUsernames(ctx context.Context, usernames []string) error
		// DeleteMultiCacheByUsernames 根据Usernames删除多条数据并清理缓存
		DeleteMultiCacheByUsernamesTx(ctx context.Context, tx *gorm_gen_dao.Query, usernames []string) error
		// DeleteMultiByUsernames 根据Usernames删除多条数据
		DeleteMultiByUsernames(ctx context.Context, usernames []string) error
		// DeleteMultiByUsernames 根据Usernames删除多条数据
		DeleteMultiByUsernamesTx(ctx context.Context, tx *gorm_gen_dao.Query, usernames []string) error
		// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
		DeleteOneCacheByID(ctx context.Context, ID string) error
		// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
		DeleteOneCacheByIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ID string) error
		// DeleteOneByID 根据ID删除一条数据
		DeleteOneByID(ctx context.Context, ID string) error
		// DeleteOneByID 根据ID删除一条数据
		DeleteOneByIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ID string) error
		// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
		DeleteMultiCacheByIDS(ctx context.Context, IDS []string) error
		// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
		DeleteMultiCacheByIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, IDS []string) error
		// DeleteMultiByIDS 根据IDS删除多条数据
		DeleteMultiByIDS(ctx context.Context, IDS []string) error
		// DeleteMultiByIDS 根据IDS删除多条数据
		DeleteMultiByIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, IDS []string) error
		// DeleteUniqueIndexCache 删除唯一索引存在的缓存
		DeleteUniqueIndexCache(ctx context.Context, data []*gorm_gen_model.SysAdmin) error
	}
	ISysAdminCache interface {
		Key(fields ...any) string
		Fetch(ctx context.Context, key string, fn func() (string, error)) (string, error)
		FetchBatch(ctx context.Context, keys []string, fn func(miss []string) (map[string]string, error)) (map[string]string, error)
		Del(ctx context.Context, key string) error
		DelBatch(ctx context.Context, keys []string) error
	}
	SysAdminRepo struct {
		db    *gorm.DB
		cache ISysAdminCache
	}
)

func NewSysAdminRepo(db *gorm.DB, cache ISysAdminCache) *SysAdminRepo {
	return &SysAdminRepo{
		db:    db,
		cache: cache,
	}
}

// CreateOne 创建一条数据
func (s *SysAdminRepo) CreateOne(ctx context.Context, data *gorm_gen_model.SysAdmin) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}

// CreateOneByTx 创建一条数据(事务)
func (s *SysAdminRepo) CreateOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.SysAdmin) error {
	dao := tx.SysAdmin
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}

// CreateBatch 批量创建数据
func (s *SysAdminRepo) CreateBatch(ctx context.Context, data []*gorm_gen_model.SysAdmin, batchSize int) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	err := dao.WithContext(ctx).CreateInBatches(data, batchSize)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOne 更新一条数据
func (s *SysAdminRepo) UpdateOne(ctx context.Context, data *gorm_gen_model.SysAdmin) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	err = s.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.SysAdmin{data})
	if err != nil {
		return err
	}
	return nil
}

// UpdateOneByTx 更新一条数据(事务)
func (s *SysAdminRepo) UpdateOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.SysAdmin) error {
	dao := tx.SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	err = s.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.SysAdmin{data})
	if err != nil {
		return err
	}
	return err
}

// DeleteOneCacheByUsername 根据username删除一条数据并清理缓存
func (s *SysAdminRepo) DeleteOneCacheByUsername(ctx context.Context, username string) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	first, err := dao.WithContext(ctx).Where(dao.Username.Eq(username)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if first == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.Username.Eq(username)).Delete()
	if err != nil {
		return err
	}
	err = s.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.SysAdmin{first})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByUsername 根据username删除一条数据并清理缓存
func (s *SysAdminRepo) DeleteOneCacheByUsernameTx(ctx context.Context, tx *gorm_gen_dao.Query, username string) error {
	dao := tx.SysAdmin
	first, err := dao.WithContext(ctx).Where(dao.Username.Eq(username)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if first == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.Username.Eq(username)).Delete()
	if err != nil {
		return err
	}
	err = s.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.SysAdmin{first})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByUsername 根据username删除一条数据
func (s *SysAdminRepo) DeleteOneByUsername(ctx context.Context, username string) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.Username.Eq(username)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByUsername 根据username删除一条数据
func (s *SysAdminRepo) DeleteOneByUsernameTx(ctx context.Context, tx *gorm_gen_dao.Query, username string) error {
	dao := tx.SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.Username.Eq(username)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByUsernames 根据usernames删除多条数据并清理缓存
func (s *SysAdminRepo) DeleteMultiCacheByUsernames(ctx context.Context, usernames []string) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	list, err := dao.WithContext(ctx).Where(dao.Username.In(usernames...)).Find()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.Username.In(usernames...)).Delete()
	if err != nil {
		return err
	}
	err = s.DeleteUniqueIndexCache(ctx, list)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByUsernames 根据usernames删除多条数据并清理缓存
func (s *SysAdminRepo) DeleteMultiCacheByUsernamesTx(ctx context.Context, tx *gorm_gen_dao.Query, usernames []string) error {
	dao := tx.SysAdmin
	list, err := dao.WithContext(ctx).Where(dao.Username.In(usernames...)).Find()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.Username.In(usernames...)).Delete()
	if err != nil {
		return err
	}
	err = s.DeleteUniqueIndexCache(ctx, list)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiByUsernames 根据usernames删除多条数据
func (s *SysAdminRepo) DeleteMultiByUsernames(ctx context.Context, usernames []string) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.Username.In(usernames...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiByUsernames 根据usernames删除多条数据
func (s *SysAdminRepo) DeleteMultiByUsernamesTx(ctx context.Context, tx *gorm_gen_dao.Query, usernames []string) error {
	dao := tx.SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.Username.In(usernames...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
func (s *SysAdminRepo) DeleteOneCacheByID(ctx context.Context, ID string) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
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
	err = s.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.SysAdmin{first})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
func (s *SysAdminRepo) DeleteOneCacheByIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ID string) error {
	dao := tx.SysAdmin
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
	err = s.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.SysAdmin{first})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByID 根据ID删除一条数据
func (s *SysAdminRepo) DeleteOneByID(ctx context.Context, ID string) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByID 根据ID删除一条数据
func (s *SysAdminRepo) DeleteOneByIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ID string) error {
	dao := tx.SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
func (s *SysAdminRepo) DeleteMultiCacheByIDS(ctx context.Context, IDS []string) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
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
	err = s.DeleteUniqueIndexCache(ctx, list)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
func (s *SysAdminRepo) DeleteMultiCacheByIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, IDS []string) error {
	dao := tx.SysAdmin
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
	err = s.DeleteUniqueIndexCache(ctx, list)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiByIDS 根据IDS删除多条数据
func (s *SysAdminRepo) DeleteMultiByIDS(ctx context.Context, IDS []string) error {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiByIDS 根据IDS删除多条数据
func (s *SysAdminRepo) DeleteMultiByIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, IDS []string) error {
	dao := tx.SysAdmin
	_, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteUniqueIndexCache 删除唯一索引存在的缓存
func (s *SysAdminRepo) DeleteUniqueIndexCache(ctx context.Context, data []*gorm_gen_model.SysAdmin) error {
	keys := make([]string, 0)
	for _, v := range data {
		keys = append(keys, s.cache.Key(cacheSysAdminByUsernamePrefix, v.Username))
		keys = append(keys, s.cache.Key(cacheSysAdminByIDPrefix, v.ID))

	}
	err := s.cache.DelBatch(ctx, keys)
	if err != nil {
		return err
	}
	return nil
}

// FindOneCacheByUsername 根据username查询一条数据并设置缓存
func (s *SysAdminRepo) FindOneCacheByUsername(ctx context.Context, username string) (*gorm_gen_model.SysAdmin, error) {
	resp := new(gorm_gen_model.SysAdmin)
	key := s.cache.Key(cacheSysAdminByUsernamePrefix, username)
	cacheValue, err := s.cache.Fetch(ctx, key, func() (string, error) {
		dao := gorm_gen_dao.Use(s.db).SysAdmin
		result, err := dao.WithContext(ctx).Where(dao.Username.Eq(username)).First()
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

// FindOneByUsername 根据username查询一条数据
func (s *SysAdminRepo) FindOneByUsername(ctx context.Context, username string) (*gorm_gen_model.SysAdmin, error) {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	result, err := dao.WithContext(ctx).Where(dao.Username.Eq(username)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return result, nil
}

// FindMultiCacheByUsernames 根据usernames查询多条数据并设置缓存
func (s *SysAdminRepo) FindMultiCacheByUsernames(ctx context.Context, usernames []string) ([]*gorm_gen_model.SysAdmin, error) {
	resp := make([]*gorm_gen_model.SysAdmin, 0)
	keys := make([]string, 0)
	keyToParam := make(map[string]string)
	for _, v := range usernames {
		key := s.cache.Key(cacheSysAdminByUsernamePrefix, v)
		keys = append(keys, key)
		keyToParam[key] = v
	}
	cacheValue, err := s.cache.FetchBatch(ctx, keys, func(miss []string) (map[string]string, error) {
		params := make([]string, 0)
		for _, v := range miss {
			params = append(params, keyToParam[v])
		}
		dao := gorm_gen_dao.Use(s.db).SysAdmin
		result, err := dao.WithContext(ctx).Where(dao.Username.In(params...)).Find()
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
			value[s.cache.Key(cacheSysAdminByUsernamePrefix, v.Username)] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new(gorm_gen_model.SysAdmin)
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}

// FindMultiByUsernames 根据usernames查询多条数据
func (s *SysAdminRepo) FindMultiByUsernames(ctx context.Context, usernames []string) ([]*gorm_gen_model.SysAdmin, error) {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	result, err := dao.WithContext(ctx).Where(dao.Username.In(usernames...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindOneCacheByID 根据ID查询一条数据并设置缓存
func (s *SysAdminRepo) FindOneCacheByID(ctx context.Context, ID string) (*gorm_gen_model.SysAdmin, error) {
	resp := new(gorm_gen_model.SysAdmin)
	key := s.cache.Key(cacheSysAdminByIDPrefix, ID)
	cacheValue, err := s.cache.Fetch(ctx, key, func() (string, error) {
		dao := gorm_gen_dao.Use(s.db).SysAdmin
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
func (s *SysAdminRepo) FindOneByID(ctx context.Context, ID string) (*gorm_gen_model.SysAdmin, error) {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	result, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return result, nil
}

// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
func (s *SysAdminRepo) FindMultiCacheByIDS(ctx context.Context, IDS []string) ([]*gorm_gen_model.SysAdmin, error) {
	resp := make([]*gorm_gen_model.SysAdmin, 0)
	keys := make([]string, 0)
	keyToParam := make(map[string]string)
	for _, v := range IDS {
		key := s.cache.Key(cacheSysAdminByIDPrefix, v)
		keys = append(keys, key)
		keyToParam[key] = v
	}
	cacheValue, err := s.cache.FetchBatch(ctx, keys, func(miss []string) (map[string]string, error) {
		params := make([]string, 0)
		for _, v := range miss {
			params = append(params, keyToParam[v])
		}
		dao := gorm_gen_dao.Use(s.db).SysAdmin
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
			value[s.cache.Key(cacheSysAdminByIDPrefix, v.ID)] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new(gorm_gen_model.SysAdmin)
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}

// FindMultiByIDS 根据IDS查询多条数据
func (s *SysAdminRepo) FindMultiByIDS(ctx context.Context, IDS []string) ([]*gorm_gen_model.SysAdmin, error) {
	dao := gorm_gen_dao.Use(s.db).SysAdmin
	result, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindMultiByPaginator 查询分页数据(通用)
func (s *SysAdminRepo) FindMultiByPaginator(ctx context.Context, params *orm.PaginatorParams) ([]*gorm_gen_model.SysAdmin, int64, error) {
	result := make([]*gorm_gen_model.SysAdmin, 0)
	var total int64
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, err
	}
	err = s.db.WithContext(ctx).Model(&gorm_gen_model.SysAdmin{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, total, nil
	}
	query := s.db.WithContext(ctx)
	order := params.ConvertToOrder()
	if order != "" {
		query = query.Order(order)
	}
	limit, offset := params.ConvertToPage()
	err = query.Limit(limit).Offset(offset).Where(queryStr, args...).Find(&result).Error
	if err != nil {
		return nil, 0, err
	}
	return result, total, err
}
