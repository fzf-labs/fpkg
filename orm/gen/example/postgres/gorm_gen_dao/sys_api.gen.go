// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gorm_gen_dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/fzf-labs/fpkg/orm/gen/example/postgres/gorm_gen_model"
)

func newSysAPI(db *gorm.DB, opts ...gen.DOOption) sysAPI {
	_sysAPI := sysAPI{}

	_sysAPI.sysAPIDo.UseDB(db, opts...)
	_sysAPI.sysAPIDo.UseModel(&gorm_gen_model.SysAPI{})

	tableName := _sysAPI.sysAPIDo.TableName()
	_sysAPI.ALL = field.NewAsterisk(tableName)
	_sysAPI.ID = field.NewString(tableName, "id")
	_sysAPI.PermissionID = field.NewString(tableName, "permission_id")
	_sysAPI.Method = field.NewString(tableName, "method")
	_sysAPI.Path = field.NewString(tableName, "path")
	_sysAPI.Desc = field.NewString(tableName, "desc")
	_sysAPI.CreatedAt = field.NewTime(tableName, "created_at")
	_sysAPI.UpdatedAt = field.NewTime(tableName, "updated_at")
	_sysAPI.DeletedAt = field.NewField(tableName, "deleted_at")

	_sysAPI.fillFieldMap()

	return _sysAPI
}

type sysAPI struct {
	sysAPIDo sysAPIDo

	ALL          field.Asterisk
	ID           field.String // 编号
	PermissionID field.String // 权限Id
	Method       field.String // 方法
	Path         field.String // 路径
	Desc         field.String // 描述
	CreatedAt    field.Time   // 创建时间
	UpdatedAt    field.Time   // 更新时间
	DeletedAt    field.Field  // 删除时间

	fieldMap map[string]field.Expr
}

func (s sysAPI) Table(newTableName string) *sysAPI {
	s.sysAPIDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s sysAPI) As(alias string) *sysAPI {
	s.sysAPIDo.DO = *(s.sysAPIDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *sysAPI) updateTableName(table string) *sysAPI {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewString(table, "id")
	s.PermissionID = field.NewString(table, "permission_id")
	s.Method = field.NewString(table, "method")
	s.Path = field.NewString(table, "path")
	s.Desc = field.NewString(table, "desc")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")
	s.DeletedAt = field.NewField(table, "deleted_at")

	s.fillFieldMap()

	return s
}

func (s *sysAPI) WithContext(ctx context.Context) *sysAPIDo { return s.sysAPIDo.WithContext(ctx) }

func (s sysAPI) TableName() string { return s.sysAPIDo.TableName() }

func (s sysAPI) Alias() string { return s.sysAPIDo.Alias() }

func (s sysAPI) Columns(cols ...field.Expr) gen.Columns { return s.sysAPIDo.Columns(cols...) }

func (s *sysAPI) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *sysAPI) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 8)
	s.fieldMap["id"] = s.ID
	s.fieldMap["permission_id"] = s.PermissionID
	s.fieldMap["method"] = s.Method
	s.fieldMap["path"] = s.Path
	s.fieldMap["desc"] = s.Desc
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
	s.fieldMap["deleted_at"] = s.DeletedAt
}

func (s sysAPI) clone(db *gorm.DB) sysAPI {
	s.sysAPIDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s sysAPI) replaceDB(db *gorm.DB) sysAPI {
	s.sysAPIDo.ReplaceDB(db)
	return s
}

type sysAPIDo struct{ gen.DO }

func (s sysAPIDo) Debug() *sysAPIDo {
	return s.withDO(s.DO.Debug())
}

func (s sysAPIDo) WithContext(ctx context.Context) *sysAPIDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s sysAPIDo) ReadDB() *sysAPIDo {
	return s.Clauses(dbresolver.Read)
}

func (s sysAPIDo) WriteDB() *sysAPIDo {
	return s.Clauses(dbresolver.Write)
}

func (s sysAPIDo) Session(config *gorm.Session) *sysAPIDo {
	return s.withDO(s.DO.Session(config))
}

func (s sysAPIDo) Clauses(conds ...clause.Expression) *sysAPIDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s sysAPIDo) Returning(value interface{}, columns ...string) *sysAPIDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s sysAPIDo) Not(conds ...gen.Condition) *sysAPIDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s sysAPIDo) Or(conds ...gen.Condition) *sysAPIDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s sysAPIDo) Select(conds ...field.Expr) *sysAPIDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s sysAPIDo) Where(conds ...gen.Condition) *sysAPIDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s sysAPIDo) Order(conds ...field.Expr) *sysAPIDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s sysAPIDo) Distinct(cols ...field.Expr) *sysAPIDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s sysAPIDo) Omit(cols ...field.Expr) *sysAPIDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s sysAPIDo) Join(table schema.Tabler, on ...field.Expr) *sysAPIDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s sysAPIDo) LeftJoin(table schema.Tabler, on ...field.Expr) *sysAPIDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s sysAPIDo) RightJoin(table schema.Tabler, on ...field.Expr) *sysAPIDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s sysAPIDo) Group(cols ...field.Expr) *sysAPIDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s sysAPIDo) Having(conds ...gen.Condition) *sysAPIDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s sysAPIDo) Limit(limit int) *sysAPIDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s sysAPIDo) Offset(offset int) *sysAPIDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s sysAPIDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *sysAPIDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s sysAPIDo) Unscoped() *sysAPIDo {
	return s.withDO(s.DO.Unscoped())
}

func (s sysAPIDo) Create(values ...*gorm_gen_model.SysAPI) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s sysAPIDo) CreateInBatches(values []*gorm_gen_model.SysAPI, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s sysAPIDo) Save(values ...*gorm_gen_model.SysAPI) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s sysAPIDo) First() (*gorm_gen_model.SysAPI, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysAPI), nil
	}
}

func (s sysAPIDo) Take() (*gorm_gen_model.SysAPI, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysAPI), nil
	}
}

func (s sysAPIDo) Last() (*gorm_gen_model.SysAPI, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysAPI), nil
	}
}

func (s sysAPIDo) Find() ([]*gorm_gen_model.SysAPI, error) {
	result, err := s.DO.Find()
	return result.([]*gorm_gen_model.SysAPI), err
}

func (s sysAPIDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*gorm_gen_model.SysAPI, err error) {
	buf := make([]*gorm_gen_model.SysAPI, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s sysAPIDo) FindInBatches(result *[]*gorm_gen_model.SysAPI, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s sysAPIDo) Attrs(attrs ...field.AssignExpr) *sysAPIDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s sysAPIDo) Assign(attrs ...field.AssignExpr) *sysAPIDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s sysAPIDo) Joins(fields ...field.RelationField) *sysAPIDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s sysAPIDo) Preload(fields ...field.RelationField) *sysAPIDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s sysAPIDo) FirstOrInit() (*gorm_gen_model.SysAPI, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysAPI), nil
	}
}

func (s sysAPIDo) FirstOrCreate() (*gorm_gen_model.SysAPI, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysAPI), nil
	}
}

func (s sysAPIDo) FindByPage(offset int, limit int) (result []*gorm_gen_model.SysAPI, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s sysAPIDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s sysAPIDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s sysAPIDo) Delete(models ...*gorm_gen_model.SysAPI) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *sysAPIDo) withDO(do gen.Dao) *sysAPIDo {
	s.DO = *do.(*gen.DO)
	return s
}