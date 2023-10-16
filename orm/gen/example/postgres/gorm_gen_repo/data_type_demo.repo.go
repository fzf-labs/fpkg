// Code generated by gen/repo. DO NOT EDIT.
// Code generated by gen/repo. DO NOT EDIT.
// Code generated by gen/repo. DO NOT EDIT.

package gorm_gen_repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/fzf-labs/fpkg/orm"
	"github.com/fzf-labs/fpkg/orm/gen/cache"
	"github.com/fzf-labs/fpkg/orm/gen/example/postgres/gorm_gen_dao"
	"github.com/fzf-labs/fpkg/orm/gen/example/postgres/gorm_gen_model"
	"gorm.io/gorm"
)

var _ IDataTypeDemoRepo = (*DataTypeDemoRepo)(nil)

var (
	cacheDataTypeDemoByIDPrefix   = "DBCache:gorm_gen:DataTypeDemoByID"
	cacheDataTypeDemoByUlIDPrefix = "DBCache:gorm_gen:DataTypeDemoByUlID"
)

type (
	IDataTypeDemoRepo interface {
		// CreateOne 创建一条数据
		CreateOne(ctx context.Context, data *gorm_gen_model.DataTypeDemo) error
		// CreateOneByTx 创建一条数据(事务)
		CreateOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.DataTypeDemo) error
		// UpsertOne Upsert一条数据
		UpsertOne(ctx context.Context, data *gorm_gen_model.DataTypeDemo) error
		// UpsertOneByTx Upsert一条数据(事务)
		UpsertOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.DataTypeDemo) error
		// CreateBatch 批量创建数据
		CreateBatch(ctx context.Context, data []*gorm_gen_model.DataTypeDemo, batchSize int) error
		// UpdateOne 更新一条数据
		UpdateOne(ctx context.Context, data *gorm_gen_model.DataTypeDemo) error
		// UpdateOne 更新一条数据(事务)
		UpdateOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.DataTypeDemo) error
		// UpdateOneWithZero 更新一条数据,包含零值
		UpdateOneWithZero(ctx context.Context, data *gorm_gen_model.DataTypeDemo) error
		// UpdateOneWithZero 更新一条数据,包含零值(事务)
		UpdateOneWithZeroByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.DataTypeDemo) error
		// FindOneCacheByID 根据ID查询一条数据并设置缓存
		FindOneCacheByID(ctx context.Context, ID string) (*gorm_gen_model.DataTypeDemo, error)
		// FindOneByID 根据ID查询一条数据
		FindOneByID(ctx context.Context, ID string) (*gorm_gen_model.DataTypeDemo, error)
		// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
		FindMultiCacheByIDS(ctx context.Context, IDS []string) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByIDS 根据IDS查询多条数据
		FindMultiByIDS(ctx context.Context, IDS []string) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByDataTypeTime 根据dataTypeTime查询多条数据
		FindMultiByDataTypeTime(ctx context.Context, dataTypeTime time.Time) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByDataTypeTimes 根据dataTypeTimes查询多条数据
		FindMultiByDataTypeTimes(ctx context.Context, dataTypeTimes []time.Time) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByCacheKey 根据_cacheKey查询多条数据
		FindMultiByCacheKey(ctx context.Context, _cacheKey string) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByCacheKeys 根据_cacheKeys查询多条数据
		FindMultiByCacheKeys(ctx context.Context, _cacheKeys []string) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByBatchAPI 根据batchAPI查询多条数据
		FindMultiByBatchAPI(ctx context.Context, batchAPI string) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByBatchAPIS 根据batchAPIS查询多条数据
		FindMultiByBatchAPIS(ctx context.Context, batchAPIS []string) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindOneCacheByUlID 根据ulID查询一条数据并设置缓存
		FindOneCacheByUlID(ctx context.Context, ulID string) (*gorm_gen_model.DataTypeDemo, error)
		// FindOneByUlID 根据ulID查询一条数据
		FindOneByUlID(ctx context.Context, ulID string) (*gorm_gen_model.DataTypeDemo, error)
		// FindMultiCacheByUlIDS 根据ulIDS查询多条数据并设置缓存
		FindMultiCacheByUlIDS(ctx context.Context, ulIDS []string) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByUlIDS 根据ulIDS查询多条数据
		FindMultiByUlIDS(ctx context.Context, ulIDS []string) ([]*gorm_gen_model.DataTypeDemo, error)
		// FindMultiByPaginator 查询分页数据(通用)
		FindMultiByPaginator(ctx context.Context, paginatorReq *orm.PaginatorReq) ([]*gorm_gen_model.DataTypeDemo, *orm.PaginatorReply, error)
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
		// DeleteOneCacheByUlID 根据ulID删除一条数据并清理缓存
		DeleteOneCacheByUlID(ctx context.Context, ulID string) error
		// DeleteOneCacheByUlID 根据ulID删除一条数据并清理缓存
		DeleteOneCacheByUlIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ulID string) error
		// DeleteOneByUlID 根据ulID删除一条数据
		DeleteOneByUlID(ctx context.Context, ulID string) error
		// DeleteOneByUlID 根据ulID删除一条数据
		DeleteOneByUlIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ulID string) error
		// DeleteMultiCacheByUlIDS 根据UlIDS删除多条数据并清理缓存
		DeleteMultiCacheByUlIDS(ctx context.Context, ulIDS []string) error
		// DeleteMultiCacheByUlIDS 根据UlIDS删除多条数据并清理缓存
		DeleteMultiCacheByUlIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, ulIDS []string) error
		// DeleteMultiByUlIDS 根据UlIDS删除多条数据
		DeleteMultiByUlIDS(ctx context.Context, ulIDS []string) error
		// DeleteMultiByUlIDS 根据UlIDS删除多条数据
		DeleteMultiByUlIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, ulIDS []string) error
		// DeleteUniqueIndexCache 删除唯一索引存在的缓存
		DeleteUniqueIndexCache(ctx context.Context, data []*gorm_gen_model.DataTypeDemo) error
	}
	DataTypeDemoRepo struct {
		db    *gorm.DB
		cache cache.IDBCache
	}
)

func NewDataTypeDemoRepo(db *gorm.DB, cache cache.IDBCache) *DataTypeDemoRepo {
	return &DataTypeDemoRepo{
		db:    db,
		cache: cache,
	}
}

// CreateOne 创建一条数据
func (d *DataTypeDemoRepo) CreateOne(ctx context.Context, data *gorm_gen_model.DataTypeDemo) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}

// CreateOneByTx 创建一条数据(事务)
func (d *DataTypeDemoRepo) CreateOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.DataTypeDemo) error {
	dao := tx.DataTypeDemo
	err := dao.WithContext(ctx).Create(data)
	if err != nil {
		return err
	}
	return nil
}

// UpsertOne Upsert一条数据
func (d *DataTypeDemoRepo) UpsertOne(ctx context.Context, data *gorm_gen_model.DataTypeDemo) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	err := dao.WithContext(ctx).Save(data)
	if err != nil {
		return err
	}
	return nil
}

// UpsertOneByTx Upsert一条数据(事务)
func (d *DataTypeDemoRepo) UpsertOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.DataTypeDemo) error {
	dao := tx.DataTypeDemo
	err := dao.WithContext(ctx).Save(data)
	if err != nil {
		return err
	}
	return nil
}

// CreateBatch 批量创建数据
func (d *DataTypeDemoRepo) CreateBatch(ctx context.Context, data []*gorm_gen_model.DataTypeDemo, batchSize int) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	err := dao.WithContext(ctx).CreateInBatches(data, batchSize)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOne 更新一条数据
func (d *DataTypeDemoRepo) UpdateOne(ctx context.Context, data *gorm_gen_model.DataTypeDemo) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.DataTypeDemo{data})
	if err != nil {
		return err
	}
	return nil
}

// UpdateOneByTx 更新一条数据(事务)
func (d *DataTypeDemoRepo) UpdateOneByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.DataTypeDemo) error {
	dao := tx.DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Updates(data)
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.DataTypeDemo{data})
	if err != nil {
		return err
	}
	return err
}

// UpdateOneWithZero 更新一条数据,包含零值
func (d *DataTypeDemoRepo) UpdateOneWithZero(ctx context.Context, data *gorm_gen_model.DataTypeDemo) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Select(dao.ALL.WithTable("")).Updates(data)
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.DataTypeDemo{data})
	if err != nil {
		return err
	}
	return nil
}

// UpdateOneWithZeroByTx 更新一条数据(事务),包含零值
func (d *DataTypeDemoRepo) UpdateOneWithZeroByTx(ctx context.Context, tx *gorm_gen_dao.Query, data *gorm_gen_model.DataTypeDemo) error {
	dao := tx.DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(data.ID)).Select(dao.ALL.WithTable("")).Updates(data)
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.DataTypeDemo{data})
	if err != nil {
		return err
	}
	return err
}

// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
func (d *DataTypeDemoRepo) DeleteOneCacheByID(ctx context.Context, ID string) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if result == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.DataTypeDemo{result})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByID 根据ID删除一条数据并清理缓存
func (d *DataTypeDemoRepo) DeleteOneCacheByIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ID string) error {
	dao := tx.DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if result == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.DataTypeDemo{result})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByID 根据ID删除一条数据
func (d *DataTypeDemoRepo) DeleteOneByID(ctx context.Context, ID string) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByID 根据ID删除一条数据
func (d *DataTypeDemoRepo) DeleteOneByIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ID string) error {
	dao := tx.DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
func (d *DataTypeDemoRepo) DeleteMultiCacheByIDS(ctx context.Context, IDS []string) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Find()
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, result)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByIDS 根据IDS删除多条数据并清理缓存
func (d *DataTypeDemoRepo) DeleteMultiCacheByIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, IDS []string) error {
	dao := tx.DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Find()
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, result)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiByIDS 根据IDS删除多条数据
func (d *DataTypeDemoRepo) DeleteMultiByIDS(ctx context.Context, IDS []string) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiByIDS 根据IDS删除多条数据
func (d *DataTypeDemoRepo) DeleteMultiByIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, IDS []string) error {
	dao := tx.DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByUlID 根据ulID删除一条数据并清理缓存
func (d *DataTypeDemoRepo) DeleteOneCacheByUlID(ctx context.Context, ulID string) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.UlID.Eq(ulID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if result == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.UlID.Eq(ulID)).Delete()
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.DataTypeDemo{result})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneCacheByUlID 根据ulID删除一条数据并清理缓存
func (d *DataTypeDemoRepo) DeleteOneCacheByUlIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ulID string) error {
	dao := tx.DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.UlID.Eq(ulID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if result == nil {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.UlID.Eq(ulID)).Delete()
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, []*gorm_gen_model.DataTypeDemo{result})
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByUlID 根据ulID删除一条数据
func (d *DataTypeDemoRepo) DeleteOneByUlID(ctx context.Context, ulID string) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.UlID.Eq(ulID)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteOneByUlID 根据ulID删除一条数据
func (d *DataTypeDemoRepo) DeleteOneByUlIDTx(ctx context.Context, tx *gorm_gen_dao.Query, ulID string) error {
	dao := tx.DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.UlID.Eq(ulID)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByUlIDS 根据ulIDS删除多条数据并清理缓存
func (d *DataTypeDemoRepo) DeleteMultiCacheByUlIDS(ctx context.Context, ulIDS []string) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.UlID.In(ulIDS...)).Find()
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.UlID.In(ulIDS...)).Delete()
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, result)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiCacheByUlIDS 根据ulIDS删除多条数据并清理缓存
func (d *DataTypeDemoRepo) DeleteMultiCacheByUlIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, ulIDS []string) error {
	dao := tx.DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.UlID.In(ulIDS...)).Find()
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return nil
	}
	_, err = dao.WithContext(ctx).Where(dao.UlID.In(ulIDS...)).Delete()
	if err != nil {
		return err
	}
	err = d.DeleteUniqueIndexCache(ctx, result)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiByUlIDS 根据ulIDS删除多条数据
func (d *DataTypeDemoRepo) DeleteMultiByUlIDS(ctx context.Context, ulIDS []string) error {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.UlID.In(ulIDS...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteMultiByUlIDS 根据ulIDS删除多条数据
func (d *DataTypeDemoRepo) DeleteMultiByUlIDSTx(ctx context.Context, tx *gorm_gen_dao.Query, ulIDS []string) error {
	dao := tx.DataTypeDemo
	_, err := dao.WithContext(ctx).Where(dao.UlID.In(ulIDS...)).Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteUniqueIndexCache 删除唯一索引存在的缓存
func (d *DataTypeDemoRepo) DeleteUniqueIndexCache(ctx context.Context, data []*gorm_gen_model.DataTypeDemo) error {
	keys := make([]string, 0)
	for _, v := range data {
		keys = append(keys, d.cache.Key(cacheDataTypeDemoByIDPrefix, v.ID))
		keys = append(keys, d.cache.Key(cacheDataTypeDemoByUlIDPrefix, v.UlID))

	}
	err := d.cache.DelBatch(ctx, keys)
	if err != nil {
		return err
	}
	return nil
}

// FindOneCacheByID 根据ID查询一条数据并设置缓存
func (d *DataTypeDemoRepo) FindOneCacheByID(ctx context.Context, ID string) (*gorm_gen_model.DataTypeDemo, error) {
	resp := new(gorm_gen_model.DataTypeDemo)
	cacheKey := d.cache.Key(cacheDataTypeDemoByIDPrefix, ID)
	cacheValue, err := d.cache.Fetch(ctx, cacheKey, func() (string, error) {
		dao := gorm_gen_dao.Use(d.db).DataTypeDemo
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
func (d *DataTypeDemoRepo) FindOneByID(ctx context.Context, ID string) (*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.ID.Eq(ID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return result, nil
}

// FindMultiCacheByIDS 根据IDS查询多条数据并设置缓存
func (d *DataTypeDemoRepo) FindMultiCacheByIDS(ctx context.Context, IDS []string) ([]*gorm_gen_model.DataTypeDemo, error) {
	resp := make([]*gorm_gen_model.DataTypeDemo, 0)
	cacheKeys := make([]string, 0)
	keyToParam := make(map[string]string)
	for _, v := range IDS {
		cacheKey := d.cache.Key(cacheDataTypeDemoByIDPrefix, v)
		cacheKeys = append(cacheKeys, cacheKey)
		keyToParam[cacheKey] = v
	}
	cacheValue, err := d.cache.FetchBatch(ctx, cacheKeys, func(miss []string) (map[string]string, error) {
		parameters := make([]string, 0)
		for _, v := range miss {
			parameters = append(parameters, keyToParam[v])
		}
		dao := gorm_gen_dao.Use(d.db).DataTypeDemo
		result, err := dao.WithContext(ctx).Where(dao.ID.In(parameters...)).Find()
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
			value[d.cache.Key(cacheDataTypeDemoByIDPrefix, v.ID)] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new(gorm_gen_model.DataTypeDemo)
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}

// FindMultiByIDS 根据IDS查询多条数据
func (d *DataTypeDemoRepo) FindMultiByIDS(ctx context.Context, IDS []string) ([]*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.ID.In(IDS...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindMultiByDataTypeTime 根据dataTypeTime查询多条数据
func (d *DataTypeDemoRepo) FindMultiByDataTypeTime(ctx context.Context, dataTypeTime time.Time) ([]*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.DataTypeTime.Eq(dataTypeTime)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindMultiByDataTypeTimes 根据dataTypeTimes查询多条数据
func (d *DataTypeDemoRepo) FindMultiByDataTypeTimes(ctx context.Context, dataTypeTimes []time.Time) ([]*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.DataTypeTime.In(dataTypeTimes...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindMultiByCacheKey 根据_cacheKey查询多条数据
func (d *DataTypeDemoRepo) FindMultiByCacheKey(ctx context.Context, _cacheKey string) ([]*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.CacheKey.Eq(_cacheKey)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindMultiByCacheKeys 根据_cacheKeys查询多条数据
func (d *DataTypeDemoRepo) FindMultiByCacheKeys(ctx context.Context, _cacheKeys []string) ([]*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.CacheKey.In(_cacheKeys...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindMultiByBatchAPI 根据batchAPI查询多条数据
func (d *DataTypeDemoRepo) FindMultiByBatchAPI(ctx context.Context, batchAPI string) ([]*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.BatchAPI.Eq(batchAPI)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindMultiByBatchAPIS 根据batchAPIS查询多条数据
func (d *DataTypeDemoRepo) FindMultiByBatchAPIS(ctx context.Context, batchAPIS []string) ([]*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.BatchAPI.In(batchAPIS...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindOneCacheByUlID 根据ulID查询一条数据并设置缓存
func (d *DataTypeDemoRepo) FindOneCacheByUlID(ctx context.Context, ulID string) (*gorm_gen_model.DataTypeDemo, error) {
	resp := new(gorm_gen_model.DataTypeDemo)
	cacheKey := d.cache.Key(cacheDataTypeDemoByUlIDPrefix, ulID)
	cacheValue, err := d.cache.Fetch(ctx, cacheKey, func() (string, error) {
		dao := gorm_gen_dao.Use(d.db).DataTypeDemo
		result, err := dao.WithContext(ctx).Where(dao.UlID.Eq(ulID)).First()
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

// FindOneByUlID 根据ulID查询一条数据
func (d *DataTypeDemoRepo) FindOneByUlID(ctx context.Context, ulID string) (*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.UlID.Eq(ulID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return result, nil
}

// FindMultiCacheByUlIDS 根据ulIDS查询多条数据并设置缓存
func (d *DataTypeDemoRepo) FindMultiCacheByUlIDS(ctx context.Context, ulIDS []string) ([]*gorm_gen_model.DataTypeDemo, error) {
	resp := make([]*gorm_gen_model.DataTypeDemo, 0)
	cacheKeys := make([]string, 0)
	keyToParam := make(map[string]string)
	for _, v := range ulIDS {
		cacheKey := d.cache.Key(cacheDataTypeDemoByUlIDPrefix, v)
		cacheKeys = append(cacheKeys, cacheKey)
		keyToParam[cacheKey] = v
	}
	cacheValue, err := d.cache.FetchBatch(ctx, cacheKeys, func(miss []string) (map[string]string, error) {
		parameters := make([]string, 0)
		for _, v := range miss {
			parameters = append(parameters, keyToParam[v])
		}
		dao := gorm_gen_dao.Use(d.db).DataTypeDemo
		result, err := dao.WithContext(ctx).Where(dao.UlID.In(parameters...)).Find()
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
			value[d.cache.Key(cacheDataTypeDemoByUlIDPrefix, v.UlID)] = string(marshal)
		}
		return value, nil
	})
	if err != nil {
		return nil, err
	}
	for _, v := range cacheValue {
		tmp := new(gorm_gen_model.DataTypeDemo)
		err := json.Unmarshal([]byte(v), tmp)
		if err != nil {
			return nil, err
		}
		resp = append(resp, tmp)
	}
	return resp, nil
}

// FindMultiByUlIDS 根据ulIDS查询多条数据
func (d *DataTypeDemoRepo) FindMultiByUlIDS(ctx context.Context, ulIDS []string) ([]*gorm_gen_model.DataTypeDemo, error) {
	dao := gorm_gen_dao.Use(d.db).DataTypeDemo
	result, err := dao.WithContext(ctx).Where(dao.UlID.In(ulIDS...)).Find()
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FindMultiByPaginator 查询分页数据(通用)
func (d *DataTypeDemoRepo) FindMultiByPaginator(ctx context.Context, paginatorReq *orm.PaginatorReq) ([]*gorm_gen_model.DataTypeDemo, *orm.PaginatorReply, error) {
	result := make([]*gorm_gen_model.DataTypeDemo, 0)
	var total int64
	queryStr, args, err := paginatorReq.ConvertToGormConditions()
	if err != nil {
		return result, nil, err
	}
	err = d.db.WithContext(ctx).Model(&gorm_gen_model.DataTypeDemo{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
	if err != nil {
		return result, nil, err
	}
	if total == 0 {
		return result, nil, nil
	}
	query := d.db.WithContext(ctx)
	order := paginatorReq.ConvertToOrder()
	if order != "" {
		query = query.Order(order)
	}
	paginatorReply := paginatorReq.ConvertToPage(int(total))
	err = query.Limit(paginatorReply.Limit).Offset(paginatorReply.Offset).Where(queryStr, args...).Find(&result).Error
	if err != nil {
		return result, nil, err
	}
	return result, paginatorReply, err
}
