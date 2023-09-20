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

func newDataTypeDemo(db *gorm.DB, opts ...gen.DOOption) dataTypeDemo {
	_dataTypeDemo := dataTypeDemo{}

	_dataTypeDemo.dataTypeDemoDo.UseDB(db, opts...)
	_dataTypeDemo.dataTypeDemoDo.UseModel(&gorm_gen_model.DataTypeDemo{})

	tableName := _dataTypeDemo.dataTypeDemoDo.TableName()
	_dataTypeDemo.ALL = field.NewAsterisk(tableName)
	_dataTypeDemo.ID = field.NewString(tableName, "id")
	_dataTypeDemo.DataTypeBool = field.NewBool(tableName, "data_type_bool")
	_dataTypeDemo.DataTypeInt2 = field.NewInt16(tableName, "data_type_int2")
	_dataTypeDemo.DataTypeInt8 = field.NewInt64(tableName, "data_type_int8")
	_dataTypeDemo.DataTypeVarchar = field.NewString(tableName, "data_type_varchar")
	_dataTypeDemo.DataTypeText = field.NewString(tableName, "data_type_text")
	_dataTypeDemo.DataTypeJSON = field.NewField(tableName, "data_type_json")
	_dataTypeDemo.CreatedAt = field.NewTime(tableName, "created_at")
	_dataTypeDemo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_dataTypeDemo.DeletedAt = field.NewField(tableName, "deleted_at")
	_dataTypeDemo.DataTypeTimeNull = field.NewField(tableName, "data_type_time_null")
	_dataTypeDemo.DataTypeTime = field.NewTime(tableName, "data_type_time")
	_dataTypeDemo.DataTypeJsonb = field.NewField(tableName, "data_type_jsonb")
	_dataTypeDemo.DataTypeDate = field.NewTime(tableName, "data_type_date")
	_dataTypeDemo.DataTypeFloat4 = field.NewFloat32(tableName, "data_type_float4")
	_dataTypeDemo.DataTypeFloat8 = field.NewFloat64(tableName, "data_type_float8")
	_dataTypeDemo.ULid = field.NewString(tableName, "_id")
	_dataTypeDemo.CacheKey = field.NewString(tableName, "cacheKey")
	_dataTypeDemo.DataTypeTimestamp = field.NewTime(tableName, "data_type_timestamp")
	_dataTypeDemo.DataTypeBytea = field.NewField(tableName, "data_type_bytea")
	_dataTypeDemo.DataTypeNumeric = field.NewFloat64(tableName, "data_type_numeric")
	_dataTypeDemo.DataTypeInterval = field.NewString(tableName, "data_type_interval")

	_dataTypeDemo.fillFieldMap()

	return _dataTypeDemo
}

type dataTypeDemo struct {
	dataTypeDemoDo dataTypeDemoDo

	ALL               field.Asterisk
	ID                field.String // ID
	DataTypeBool      field.Bool   // 数据类型 bool
	DataTypeInt2      field.Int16  // 数据类型 int2
	DataTypeInt8      field.Int64  // 数据类型 int8
	DataTypeVarchar   field.String // 数据类型 varchar
	DataTypeText      field.String // 数据类型 text
	DataTypeJSON      field.Field  // 数据类型 json
	CreatedAt         field.Time   // 创建时间
	UpdatedAt         field.Time   // 更新时间
	DeletedAt         field.Field  // 删除时间
	DataTypeTimeNull  field.Field  // 数据类型 time null
	DataTypeTime      field.Time   // 数据类型 time not null
	DataTypeJsonb     field.Field  // 数据类型 jsonb
	DataTypeDate      field.Time
	DataTypeFloat4    field.Float32
	DataTypeFloat8    field.Float64
	ULid              field.String // 验证下划线
	CacheKey          field.String // 特殊保留字段名称
	DataTypeTimestamp field.Time
	DataTypeBytea     field.Field
	DataTypeNumeric   field.Float64
	DataTypeInterval  field.String

	fieldMap map[string]field.Expr
}

func (d dataTypeDemo) Table(newTableName string) *dataTypeDemo {
	d.dataTypeDemoDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d dataTypeDemo) As(alias string) *dataTypeDemo {
	d.dataTypeDemoDo.DO = *(d.dataTypeDemoDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *dataTypeDemo) updateTableName(table string) *dataTypeDemo {
	d.ALL = field.NewAsterisk(table)
	d.ID = field.NewString(table, "id")
	d.DataTypeBool = field.NewBool(table, "data_type_bool")
	d.DataTypeInt2 = field.NewInt16(table, "data_type_int2")
	d.DataTypeInt8 = field.NewInt64(table, "data_type_int8")
	d.DataTypeVarchar = field.NewString(table, "data_type_varchar")
	d.DataTypeText = field.NewString(table, "data_type_text")
	d.DataTypeJSON = field.NewField(table, "data_type_json")
	d.CreatedAt = field.NewTime(table, "created_at")
	d.UpdatedAt = field.NewTime(table, "updated_at")
	d.DeletedAt = field.NewField(table, "deleted_at")
	d.DataTypeTimeNull = field.NewField(table, "data_type_time_null")
	d.DataTypeTime = field.NewTime(table, "data_type_time")
	d.DataTypeJsonb = field.NewField(table, "data_type_jsonb")
	d.DataTypeDate = field.NewTime(table, "data_type_date")
	d.DataTypeFloat4 = field.NewFloat32(table, "data_type_float4")
	d.DataTypeFloat8 = field.NewFloat64(table, "data_type_float8")
	d.ULid = field.NewString(table, "_id")
	d.CacheKey = field.NewString(table, "cacheKey")
	d.DataTypeTimestamp = field.NewTime(table, "data_type_timestamp")
	d.DataTypeBytea = field.NewField(table, "data_type_bytea")
	d.DataTypeNumeric = field.NewFloat64(table, "data_type_numeric")
	d.DataTypeInterval = field.NewString(table, "data_type_interval")

	d.fillFieldMap()

	return d
}

func (d *dataTypeDemo) WithContext(ctx context.Context) *dataTypeDemoDo {
	return d.dataTypeDemoDo.WithContext(ctx)
}

func (d dataTypeDemo) TableName() string { return d.dataTypeDemoDo.TableName() }

func (d dataTypeDemo) Alias() string { return d.dataTypeDemoDo.Alias() }

func (d dataTypeDemo) Columns(cols ...field.Expr) gen.Columns {
	return d.dataTypeDemoDo.Columns(cols...)
}

func (d *dataTypeDemo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *dataTypeDemo) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 22)
	d.fieldMap["id"] = d.ID
	d.fieldMap["data_type_bool"] = d.DataTypeBool
	d.fieldMap["data_type_int2"] = d.DataTypeInt2
	d.fieldMap["data_type_int8"] = d.DataTypeInt8
	d.fieldMap["data_type_varchar"] = d.DataTypeVarchar
	d.fieldMap["data_type_text"] = d.DataTypeText
	d.fieldMap["data_type_json"] = d.DataTypeJSON
	d.fieldMap["created_at"] = d.CreatedAt
	d.fieldMap["updated_at"] = d.UpdatedAt
	d.fieldMap["deleted_at"] = d.DeletedAt
	d.fieldMap["data_type_time_null"] = d.DataTypeTimeNull
	d.fieldMap["data_type_time"] = d.DataTypeTime
	d.fieldMap["data_type_jsonb"] = d.DataTypeJsonb
	d.fieldMap["data_type_date"] = d.DataTypeDate
	d.fieldMap["data_type_float4"] = d.DataTypeFloat4
	d.fieldMap["data_type_float8"] = d.DataTypeFloat8
	d.fieldMap["_id"] = d.ULid
	d.fieldMap["cacheKey"] = d.CacheKey
	d.fieldMap["data_type_timestamp"] = d.DataTypeTimestamp
	d.fieldMap["data_type_bytea"] = d.DataTypeBytea
	d.fieldMap["data_type_numeric"] = d.DataTypeNumeric
	d.fieldMap["data_type_interval"] = d.DataTypeInterval
}

func (d dataTypeDemo) clone(db *gorm.DB) dataTypeDemo {
	d.dataTypeDemoDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d dataTypeDemo) replaceDB(db *gorm.DB) dataTypeDemo {
	d.dataTypeDemoDo.ReplaceDB(db)
	return d
}

type dataTypeDemoDo struct{ gen.DO }

func (d dataTypeDemoDo) Debug() *dataTypeDemoDo {
	return d.withDO(d.DO.Debug())
}

func (d dataTypeDemoDo) WithContext(ctx context.Context) *dataTypeDemoDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d dataTypeDemoDo) ReadDB() *dataTypeDemoDo {
	return d.Clauses(dbresolver.Read)
}

func (d dataTypeDemoDo) WriteDB() *dataTypeDemoDo {
	return d.Clauses(dbresolver.Write)
}

func (d dataTypeDemoDo) Session(config *gorm.Session) *dataTypeDemoDo {
	return d.withDO(d.DO.Session(config))
}

func (d dataTypeDemoDo) Clauses(conds ...clause.Expression) *dataTypeDemoDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d dataTypeDemoDo) Returning(value interface{}, columns ...string) *dataTypeDemoDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d dataTypeDemoDo) Not(conds ...gen.Condition) *dataTypeDemoDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d dataTypeDemoDo) Or(conds ...gen.Condition) *dataTypeDemoDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d dataTypeDemoDo) Select(conds ...field.Expr) *dataTypeDemoDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d dataTypeDemoDo) Where(conds ...gen.Condition) *dataTypeDemoDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d dataTypeDemoDo) Order(conds ...field.Expr) *dataTypeDemoDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d dataTypeDemoDo) Distinct(cols ...field.Expr) *dataTypeDemoDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d dataTypeDemoDo) Omit(cols ...field.Expr) *dataTypeDemoDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d dataTypeDemoDo) Join(table schema.Tabler, on ...field.Expr) *dataTypeDemoDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d dataTypeDemoDo) LeftJoin(table schema.Tabler, on ...field.Expr) *dataTypeDemoDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d dataTypeDemoDo) RightJoin(table schema.Tabler, on ...field.Expr) *dataTypeDemoDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d dataTypeDemoDo) Group(cols ...field.Expr) *dataTypeDemoDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d dataTypeDemoDo) Having(conds ...gen.Condition) *dataTypeDemoDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d dataTypeDemoDo) Limit(limit int) *dataTypeDemoDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d dataTypeDemoDo) Offset(offset int) *dataTypeDemoDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d dataTypeDemoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *dataTypeDemoDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d dataTypeDemoDo) Unscoped() *dataTypeDemoDo {
	return d.withDO(d.DO.Unscoped())
}

func (d dataTypeDemoDo) Create(values ...*gorm_gen_model.DataTypeDemo) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d dataTypeDemoDo) CreateInBatches(values []*gorm_gen_model.DataTypeDemo, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d dataTypeDemoDo) Save(values ...*gorm_gen_model.DataTypeDemo) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d dataTypeDemoDo) First() (*gorm_gen_model.DataTypeDemo, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.DataTypeDemo), nil
	}
}

func (d dataTypeDemoDo) Take() (*gorm_gen_model.DataTypeDemo, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.DataTypeDemo), nil
	}
}

func (d dataTypeDemoDo) Last() (*gorm_gen_model.DataTypeDemo, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.DataTypeDemo), nil
	}
}

func (d dataTypeDemoDo) Find() ([]*gorm_gen_model.DataTypeDemo, error) {
	result, err := d.DO.Find()
	return result.([]*gorm_gen_model.DataTypeDemo), err
}

func (d dataTypeDemoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*gorm_gen_model.DataTypeDemo, err error) {
	buf := make([]*gorm_gen_model.DataTypeDemo, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d dataTypeDemoDo) FindInBatches(result *[]*gorm_gen_model.DataTypeDemo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d dataTypeDemoDo) Attrs(attrs ...field.AssignExpr) *dataTypeDemoDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d dataTypeDemoDo) Assign(attrs ...field.AssignExpr) *dataTypeDemoDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d dataTypeDemoDo) Joins(fields ...field.RelationField) *dataTypeDemoDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d dataTypeDemoDo) Preload(fields ...field.RelationField) *dataTypeDemoDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d dataTypeDemoDo) FirstOrInit() (*gorm_gen_model.DataTypeDemo, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.DataTypeDemo), nil
	}
}

func (d dataTypeDemoDo) FirstOrCreate() (*gorm_gen_model.DataTypeDemo, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.DataTypeDemo), nil
	}
}

func (d dataTypeDemoDo) FindByPage(offset int, limit int) (result []*gorm_gen_model.DataTypeDemo, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d dataTypeDemoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d dataTypeDemoDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d dataTypeDemoDo) Delete(models ...*gorm_gen_model.DataTypeDemo) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *dataTypeDemoDo) withDO(do gen.Dao) *dataTypeDemoDo {
	d.DO = *do.(*gen.DO)
	return d
}
