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

func newSysDept(db *gorm.DB, opts ...gen.DOOption) sysDept {
	_sysDept := sysDept{}

	_sysDept.sysDeptDo.UseDB(db, opts...)
	_sysDept.sysDeptDo.UseModel(&gorm_gen_model.SysDept{})

	tableName := _sysDept.sysDeptDo.TableName()
	_sysDept.ALL = field.NewAsterisk(tableName)
	_sysDept.ID = field.NewString(tableName, "id")
	_sysDept.Pid = field.NewString(tableName, "pid")
	_sysDept.Name = field.NewString(tableName, "name")
	_sysDept.FullName = field.NewString(tableName, "full_name")
	_sysDept.Responsible = field.NewString(tableName, "responsible")
	_sysDept.Phone = field.NewString(tableName, "phone")
	_sysDept.Email = field.NewString(tableName, "email")
	_sysDept.Type = field.NewInt16(tableName, "type")
	_sysDept.Status = field.NewInt16(tableName, "status")
	_sysDept.Sort = field.NewInt64(tableName, "sort")
	_sysDept.CreatedAt = field.NewTime(tableName, "created_at")
	_sysDept.UpdatedAt = field.NewTime(tableName, "updated_at")
	_sysDept.DeletedAt = field.NewField(tableName, "deleted_at")

	_sysDept.fillFieldMap()

	return _sysDept
}

type sysDept struct {
	sysDeptDo sysDeptDo

	ALL         field.Asterisk
	ID          field.String // 编号
	Pid         field.String // 父级id
	Name        field.String // 部门简称
	FullName    field.String // 部门全称
	Responsible field.String // 负责人
	Phone       field.String // 负责人电话
	Email       field.String // 负责人邮箱
	Type        field.Int16  // 1=公司 2=子公司 3=部门
	Status      field.Int16  // 0=禁用 1=开启
	Sort        field.Int64  // 排序值
	CreatedAt   field.Time   // 创建时间
	UpdatedAt   field.Time   // 更新时间
	DeletedAt   field.Field  // 删除时间

	fieldMap map[string]field.Expr
}

func (s sysDept) Table(newTableName string) *sysDept {
	s.sysDeptDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s sysDept) As(alias string) *sysDept {
	s.sysDeptDo.DO = *(s.sysDeptDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *sysDept) updateTableName(table string) *sysDept {
	s.ALL = field.NewAsterisk(table)
	s.ID = field.NewString(table, "id")
	s.Pid = field.NewString(table, "pid")
	s.Name = field.NewString(table, "name")
	s.FullName = field.NewString(table, "full_name")
	s.Responsible = field.NewString(table, "responsible")
	s.Phone = field.NewString(table, "phone")
	s.Email = field.NewString(table, "email")
	s.Type = field.NewInt16(table, "type")
	s.Status = field.NewInt16(table, "status")
	s.Sort = field.NewInt64(table, "sort")
	s.CreatedAt = field.NewTime(table, "created_at")
	s.UpdatedAt = field.NewTime(table, "updated_at")
	s.DeletedAt = field.NewField(table, "deleted_at")

	s.fillFieldMap()

	return s
}

func (s *sysDept) WithContext(ctx context.Context) *sysDeptDo { return s.sysDeptDo.WithContext(ctx) }

func (s sysDept) TableName() string { return s.sysDeptDo.TableName() }

func (s sysDept) Alias() string { return s.sysDeptDo.Alias() }

func (s sysDept) Columns(cols ...field.Expr) gen.Columns { return s.sysDeptDo.Columns(cols...) }

func (s *sysDept) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *sysDept) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 13)
	s.fieldMap["id"] = s.ID
	s.fieldMap["pid"] = s.Pid
	s.fieldMap["name"] = s.Name
	s.fieldMap["full_name"] = s.FullName
	s.fieldMap["responsible"] = s.Responsible
	s.fieldMap["phone"] = s.Phone
	s.fieldMap["email"] = s.Email
	s.fieldMap["type"] = s.Type
	s.fieldMap["status"] = s.Status
	s.fieldMap["sort"] = s.Sort
	s.fieldMap["created_at"] = s.CreatedAt
	s.fieldMap["updated_at"] = s.UpdatedAt
	s.fieldMap["deleted_at"] = s.DeletedAt
}

func (s sysDept) clone(db *gorm.DB) sysDept {
	s.sysDeptDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s sysDept) replaceDB(db *gorm.DB) sysDept {
	s.sysDeptDo.ReplaceDB(db)
	return s
}

type sysDeptDo struct{ gen.DO }

func (s sysDeptDo) Debug() *sysDeptDo {
	return s.withDO(s.DO.Debug())
}

func (s sysDeptDo) WithContext(ctx context.Context) *sysDeptDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s sysDeptDo) ReadDB() *sysDeptDo {
	return s.Clauses(dbresolver.Read)
}

func (s sysDeptDo) WriteDB() *sysDeptDo {
	return s.Clauses(dbresolver.Write)
}

func (s sysDeptDo) Session(config *gorm.Session) *sysDeptDo {
	return s.withDO(s.DO.Session(config))
}

func (s sysDeptDo) Clauses(conds ...clause.Expression) *sysDeptDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s sysDeptDo) Returning(value interface{}, columns ...string) *sysDeptDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s sysDeptDo) Not(conds ...gen.Condition) *sysDeptDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s sysDeptDo) Or(conds ...gen.Condition) *sysDeptDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s sysDeptDo) Select(conds ...field.Expr) *sysDeptDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s sysDeptDo) Where(conds ...gen.Condition) *sysDeptDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s sysDeptDo) Order(conds ...field.Expr) *sysDeptDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s sysDeptDo) Distinct(cols ...field.Expr) *sysDeptDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s sysDeptDo) Omit(cols ...field.Expr) *sysDeptDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s sysDeptDo) Join(table schema.Tabler, on ...field.Expr) *sysDeptDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s sysDeptDo) LeftJoin(table schema.Tabler, on ...field.Expr) *sysDeptDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s sysDeptDo) RightJoin(table schema.Tabler, on ...field.Expr) *sysDeptDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s sysDeptDo) Group(cols ...field.Expr) *sysDeptDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s sysDeptDo) Having(conds ...gen.Condition) *sysDeptDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s sysDeptDo) Limit(limit int) *sysDeptDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s sysDeptDo) Offset(offset int) *sysDeptDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s sysDeptDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *sysDeptDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s sysDeptDo) Unscoped() *sysDeptDo {
	return s.withDO(s.DO.Unscoped())
}

func (s sysDeptDo) Create(values ...*gorm_gen_model.SysDept) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s sysDeptDo) CreateInBatches(values []*gorm_gen_model.SysDept, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s sysDeptDo) Save(values ...*gorm_gen_model.SysDept) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s sysDeptDo) First() (*gorm_gen_model.SysDept, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysDept), nil
	}
}

func (s sysDeptDo) Take() (*gorm_gen_model.SysDept, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysDept), nil
	}
}

func (s sysDeptDo) Last() (*gorm_gen_model.SysDept, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysDept), nil
	}
}

func (s sysDeptDo) Find() ([]*gorm_gen_model.SysDept, error) {
	result, err := s.DO.Find()
	return result.([]*gorm_gen_model.SysDept), err
}

func (s sysDeptDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*gorm_gen_model.SysDept, err error) {
	buf := make([]*gorm_gen_model.SysDept, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s sysDeptDo) FindInBatches(result *[]*gorm_gen_model.SysDept, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s sysDeptDo) Attrs(attrs ...field.AssignExpr) *sysDeptDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s sysDeptDo) Assign(attrs ...field.AssignExpr) *sysDeptDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s sysDeptDo) Joins(fields ...field.RelationField) *sysDeptDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s sysDeptDo) Preload(fields ...field.RelationField) *sysDeptDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s sysDeptDo) FirstOrInit() (*gorm_gen_model.SysDept, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysDept), nil
	}
}

func (s sysDeptDo) FirstOrCreate() (*gorm_gen_model.SysDept, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.SysDept), nil
	}
}

func (s sysDeptDo) FindByPage(offset int, limit int) (result []*gorm_gen_model.SysDept, count int64, err error) {
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

func (s sysDeptDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s sysDeptDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s sysDeptDo) Delete(models ...*gorm_gen_model.SysDept) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *sysDeptDo) withDO(do gen.Dao) *sysDeptDo {
	s.DO = *do.(*gen.DO)
	return s
}