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

func newAdminDemo(db *gorm.DB, opts ...gen.DOOption) adminDemo {
	_adminDemo := adminDemo{}

	_adminDemo.adminDemoDo.UseDB(db, opts...)
	_adminDemo.adminDemoDo.UseModel(&gorm_gen_model.AdminDemo{})

	tableName := _adminDemo.adminDemoDo.TableName()
	_adminDemo.ALL = field.NewAsterisk(tableName)
	_adminDemo.ID = field.NewString(tableName, "id")
	_adminDemo.Username = field.NewString(tableName, "username")
	_adminDemo.Password = field.NewString(tableName, "password")
	_adminDemo.Nickname = field.NewString(tableName, "nickname")
	_adminDemo.Avatar = field.NewString(tableName, "avatar")
	_adminDemo.Gender = field.NewInt16(tableName, "gender")
	_adminDemo.Email = field.NewString(tableName, "email")
	_adminDemo.Mobile = field.NewString(tableName, "mobile")
	_adminDemo.JobID = field.NewString(tableName, "job_id")
	_adminDemo.DeptID = field.NewString(tableName, "dept_id")
	_adminDemo.RoleIds = field.NewField(tableName, "role_ids")
	_adminDemo.Salt = field.NewString(tableName, "salt")
	_adminDemo.Status = field.NewInt16(tableName, "status")
	_adminDemo.Motto = field.NewString(tableName, "motto")
	_adminDemo.CreatedAt = field.NewTime(tableName, "created_at")
	_adminDemo.UpdatedAt = field.NewTime(tableName, "updated_at")
	_adminDemo.DeletedAt = field.NewField(tableName, "deleted_at")
	_adminDemo.AdminLogDemos = adminDemoHasManyAdminLogDemos{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("AdminLogDemos", "gorm_gen_model.AdminLogDemo"),
	}

	_adminDemo.AdminRoles = adminDemoManyToManyAdminRoles{
		db: db.Session(&gorm.Session{}),

		RelationField: field.NewRelation("AdminRoles", "gorm_gen_model.AdminRoleDemo"),
		Admins: struct {
			field.RelationField
		}{
			RelationField: field.NewRelation("AdminRoles.Admins", "gorm_gen_model.AdminDemo"),
		},
	}

	_adminDemo.fillFieldMap()

	return _adminDemo
}

type adminDemo struct {
	adminDemoDo adminDemoDo

	ALL           field.Asterisk
	ID            field.String // 编号
	Username      field.String // 用户名
	Password      field.String // 密码
	Nickname      field.String // 昵称
	Avatar        field.String // 头像
	Gender        field.Int16  // 0=保密 1=女 2=男
	Email         field.String // 邮件
	Mobile        field.String // 手机号
	JobID         field.String // 岗位
	DeptID        field.String // 部门
	RoleIds       field.Field  // 角色集
	Salt          field.String // 盐值
	Status        field.Int16  // 0=禁用 1=开启
	Motto         field.String // 个性签名
	CreatedAt     field.Time   // 创建时间
	UpdatedAt     field.Time   // 更新时间
	DeletedAt     field.Field  // 删除时间
	AdminLogDemos adminDemoHasManyAdminLogDemos

	AdminRoles adminDemoManyToManyAdminRoles

	fieldMap map[string]field.Expr
}

func (a adminDemo) Table(newTableName string) *adminDemo {
	a.adminDemoDo.UseTable(newTableName)
	return a.updateTableName(newTableName)
}

func (a adminDemo) As(alias string) *adminDemo {
	a.adminDemoDo.DO = *(a.adminDemoDo.As(alias).(*gen.DO))
	return a.updateTableName(alias)
}

func (a *adminDemo) updateTableName(table string) *adminDemo {
	a.ALL = field.NewAsterisk(table)
	a.ID = field.NewString(table, "id")
	a.Username = field.NewString(table, "username")
	a.Password = field.NewString(table, "password")
	a.Nickname = field.NewString(table, "nickname")
	a.Avatar = field.NewString(table, "avatar")
	a.Gender = field.NewInt16(table, "gender")
	a.Email = field.NewString(table, "email")
	a.Mobile = field.NewString(table, "mobile")
	a.JobID = field.NewString(table, "job_id")
	a.DeptID = field.NewString(table, "dept_id")
	a.RoleIds = field.NewField(table, "role_ids")
	a.Salt = field.NewString(table, "salt")
	a.Status = field.NewInt16(table, "status")
	a.Motto = field.NewString(table, "motto")
	a.CreatedAt = field.NewTime(table, "created_at")
	a.UpdatedAt = field.NewTime(table, "updated_at")
	a.DeletedAt = field.NewField(table, "deleted_at")

	a.fillFieldMap()

	return a
}

func (a *adminDemo) WithContext(ctx context.Context) *adminDemoDo {
	return a.adminDemoDo.WithContext(ctx)
}

func (a adminDemo) TableName() string { return a.adminDemoDo.TableName() }

func (a adminDemo) Alias() string { return a.adminDemoDo.Alias() }

func (a adminDemo) Columns(cols ...field.Expr) gen.Columns { return a.adminDemoDo.Columns(cols...) }

func (a *adminDemo) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := a.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (a *adminDemo) fillFieldMap() {
	a.fieldMap = make(map[string]field.Expr, 19)
	a.fieldMap["id"] = a.ID
	a.fieldMap["username"] = a.Username
	a.fieldMap["password"] = a.Password
	a.fieldMap["nickname"] = a.Nickname
	a.fieldMap["avatar"] = a.Avatar
	a.fieldMap["gender"] = a.Gender
	a.fieldMap["email"] = a.Email
	a.fieldMap["mobile"] = a.Mobile
	a.fieldMap["job_id"] = a.JobID
	a.fieldMap["dept_id"] = a.DeptID
	a.fieldMap["role_ids"] = a.RoleIds
	a.fieldMap["salt"] = a.Salt
	a.fieldMap["status"] = a.Status
	a.fieldMap["motto"] = a.Motto
	a.fieldMap["created_at"] = a.CreatedAt
	a.fieldMap["updated_at"] = a.UpdatedAt
	a.fieldMap["deleted_at"] = a.DeletedAt

}

func (a adminDemo) clone(db *gorm.DB) adminDemo {
	a.adminDemoDo.ReplaceConnPool(db.Statement.ConnPool)
	return a
}

func (a adminDemo) replaceDB(db *gorm.DB) adminDemo {
	a.adminDemoDo.ReplaceDB(db)
	return a
}

type adminDemoHasManyAdminLogDemos struct {
	db *gorm.DB

	field.RelationField
}

func (a adminDemoHasManyAdminLogDemos) Where(conds ...field.Expr) *adminDemoHasManyAdminLogDemos {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a adminDemoHasManyAdminLogDemos) WithContext(ctx context.Context) *adminDemoHasManyAdminLogDemos {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a adminDemoHasManyAdminLogDemos) Session(session *gorm.Session) *adminDemoHasManyAdminLogDemos {
	a.db = a.db.Session(session)
	return &a
}

func (a adminDemoHasManyAdminLogDemos) Model(m *gorm_gen_model.AdminDemo) *adminDemoHasManyAdminLogDemosTx {
	return &adminDemoHasManyAdminLogDemosTx{a.db.Model(m).Association(a.Name())}
}

type adminDemoHasManyAdminLogDemosTx struct{ tx *gorm.Association }

func (a adminDemoHasManyAdminLogDemosTx) Find() (result []*gorm_gen_model.AdminLogDemo, err error) {
	return result, a.tx.Find(&result)
}

func (a adminDemoHasManyAdminLogDemosTx) Append(values ...*gorm_gen_model.AdminLogDemo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a adminDemoHasManyAdminLogDemosTx) Replace(values ...*gorm_gen_model.AdminLogDemo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a adminDemoHasManyAdminLogDemosTx) Delete(values ...*gorm_gen_model.AdminLogDemo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a adminDemoHasManyAdminLogDemosTx) Clear() error {
	return a.tx.Clear()
}

func (a adminDemoHasManyAdminLogDemosTx) Count() int64 {
	return a.tx.Count()
}

type adminDemoManyToManyAdminRoles struct {
	db *gorm.DB

	field.RelationField

	Admins struct {
		field.RelationField
	}
}

func (a adminDemoManyToManyAdminRoles) Where(conds ...field.Expr) *adminDemoManyToManyAdminRoles {
	if len(conds) == 0 {
		return &a
	}

	exprs := make([]clause.Expression, 0, len(conds))
	for _, cond := range conds {
		exprs = append(exprs, cond.BeCond().(clause.Expression))
	}
	a.db = a.db.Clauses(clause.Where{Exprs: exprs})
	return &a
}

func (a adminDemoManyToManyAdminRoles) WithContext(ctx context.Context) *adminDemoManyToManyAdminRoles {
	a.db = a.db.WithContext(ctx)
	return &a
}

func (a adminDemoManyToManyAdminRoles) Session(session *gorm.Session) *adminDemoManyToManyAdminRoles {
	a.db = a.db.Session(session)
	return &a
}

func (a adminDemoManyToManyAdminRoles) Model(m *gorm_gen_model.AdminDemo) *adminDemoManyToManyAdminRolesTx {
	return &adminDemoManyToManyAdminRolesTx{a.db.Model(m).Association(a.Name())}
}

type adminDemoManyToManyAdminRolesTx struct{ tx *gorm.Association }

func (a adminDemoManyToManyAdminRolesTx) Find() (result []*gorm_gen_model.AdminRoleDemo, err error) {
	return result, a.tx.Find(&result)
}

func (a adminDemoManyToManyAdminRolesTx) Append(values ...*gorm_gen_model.AdminRoleDemo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Append(targetValues...)
}

func (a adminDemoManyToManyAdminRolesTx) Replace(values ...*gorm_gen_model.AdminRoleDemo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Replace(targetValues...)
}

func (a adminDemoManyToManyAdminRolesTx) Delete(values ...*gorm_gen_model.AdminRoleDemo) (err error) {
	targetValues := make([]interface{}, len(values))
	for i, v := range values {
		targetValues[i] = v
	}
	return a.tx.Delete(targetValues...)
}

func (a adminDemoManyToManyAdminRolesTx) Clear() error {
	return a.tx.Clear()
}

func (a adminDemoManyToManyAdminRolesTx) Count() int64 {
	return a.tx.Count()
}

type adminDemoDo struct{ gen.DO }

func (a adminDemoDo) Debug() *adminDemoDo {
	return a.withDO(a.DO.Debug())
}

func (a adminDemoDo) WithContext(ctx context.Context) *adminDemoDo {
	return a.withDO(a.DO.WithContext(ctx))
}

func (a adminDemoDo) ReadDB() *adminDemoDo {
	return a.Clauses(dbresolver.Read)
}

func (a adminDemoDo) WriteDB() *adminDemoDo {
	return a.Clauses(dbresolver.Write)
}

func (a adminDemoDo) Session(config *gorm.Session) *adminDemoDo {
	return a.withDO(a.DO.Session(config))
}

func (a adminDemoDo) Clauses(conds ...clause.Expression) *adminDemoDo {
	return a.withDO(a.DO.Clauses(conds...))
}

func (a adminDemoDo) Returning(value interface{}, columns ...string) *adminDemoDo {
	return a.withDO(a.DO.Returning(value, columns...))
}

func (a adminDemoDo) Not(conds ...gen.Condition) *adminDemoDo {
	return a.withDO(a.DO.Not(conds...))
}

func (a adminDemoDo) Or(conds ...gen.Condition) *adminDemoDo {
	return a.withDO(a.DO.Or(conds...))
}

func (a adminDemoDo) Select(conds ...field.Expr) *adminDemoDo {
	return a.withDO(a.DO.Select(conds...))
}

func (a adminDemoDo) Where(conds ...gen.Condition) *adminDemoDo {
	return a.withDO(a.DO.Where(conds...))
}

func (a adminDemoDo) Order(conds ...field.Expr) *adminDemoDo {
	return a.withDO(a.DO.Order(conds...))
}

func (a adminDemoDo) Distinct(cols ...field.Expr) *adminDemoDo {
	return a.withDO(a.DO.Distinct(cols...))
}

func (a adminDemoDo) Omit(cols ...field.Expr) *adminDemoDo {
	return a.withDO(a.DO.Omit(cols...))
}

func (a adminDemoDo) Join(table schema.Tabler, on ...field.Expr) *adminDemoDo {
	return a.withDO(a.DO.Join(table, on...))
}

func (a adminDemoDo) LeftJoin(table schema.Tabler, on ...field.Expr) *adminDemoDo {
	return a.withDO(a.DO.LeftJoin(table, on...))
}

func (a adminDemoDo) RightJoin(table schema.Tabler, on ...field.Expr) *adminDemoDo {
	return a.withDO(a.DO.RightJoin(table, on...))
}

func (a adminDemoDo) Group(cols ...field.Expr) *adminDemoDo {
	return a.withDO(a.DO.Group(cols...))
}

func (a adminDemoDo) Having(conds ...gen.Condition) *adminDemoDo {
	return a.withDO(a.DO.Having(conds...))
}

func (a adminDemoDo) Limit(limit int) *adminDemoDo {
	return a.withDO(a.DO.Limit(limit))
}

func (a adminDemoDo) Offset(offset int) *adminDemoDo {
	return a.withDO(a.DO.Offset(offset))
}

func (a adminDemoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *adminDemoDo {
	return a.withDO(a.DO.Scopes(funcs...))
}

func (a adminDemoDo) Unscoped() *adminDemoDo {
	return a.withDO(a.DO.Unscoped())
}

func (a adminDemoDo) Create(values ...*gorm_gen_model.AdminDemo) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Create(values)
}

func (a adminDemoDo) CreateInBatches(values []*gorm_gen_model.AdminDemo, batchSize int) error {
	return a.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (a adminDemoDo) Save(values ...*gorm_gen_model.AdminDemo) error {
	if len(values) == 0 {
		return nil
	}
	return a.DO.Save(values)
}

func (a adminDemoDo) First() (*gorm_gen_model.AdminDemo, error) {
	if result, err := a.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.AdminDemo), nil
	}
}

func (a adminDemoDo) Take() (*gorm_gen_model.AdminDemo, error) {
	if result, err := a.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.AdminDemo), nil
	}
}

func (a adminDemoDo) Last() (*gorm_gen_model.AdminDemo, error) {
	if result, err := a.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.AdminDemo), nil
	}
}

func (a adminDemoDo) Find() ([]*gorm_gen_model.AdminDemo, error) {
	result, err := a.DO.Find()
	return result.([]*gorm_gen_model.AdminDemo), err
}

func (a adminDemoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*gorm_gen_model.AdminDemo, err error) {
	buf := make([]*gorm_gen_model.AdminDemo, 0, batchSize)
	err = a.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (a adminDemoDo) FindInBatches(result *[]*gorm_gen_model.AdminDemo, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return a.DO.FindInBatches(result, batchSize, fc)
}

func (a adminDemoDo) Attrs(attrs ...field.AssignExpr) *adminDemoDo {
	return a.withDO(a.DO.Attrs(attrs...))
}

func (a adminDemoDo) Assign(attrs ...field.AssignExpr) *adminDemoDo {
	return a.withDO(a.DO.Assign(attrs...))
}

func (a adminDemoDo) Joins(fields ...field.RelationField) *adminDemoDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Joins(_f))
	}
	return &a
}

func (a adminDemoDo) Preload(fields ...field.RelationField) *adminDemoDo {
	for _, _f := range fields {
		a = *a.withDO(a.DO.Preload(_f))
	}
	return &a
}

func (a adminDemoDo) FirstOrInit() (*gorm_gen_model.AdminDemo, error) {
	if result, err := a.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.AdminDemo), nil
	}
}

func (a adminDemoDo) FirstOrCreate() (*gorm_gen_model.AdminDemo, error) {
	if result, err := a.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*gorm_gen_model.AdminDemo), nil
	}
}

func (a adminDemoDo) FindByPage(offset int, limit int) (result []*gorm_gen_model.AdminDemo, count int64, err error) {
	result, err = a.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = a.Offset(-1).Limit(-1).Count()
	return
}

func (a adminDemoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = a.Count()
	if err != nil {
		return
	}

	err = a.Offset(offset).Limit(limit).Scan(result)
	return
}

func (a adminDemoDo) Scan(result interface{}) (err error) {
	return a.DO.Scan(result)
}

func (a adminDemoDo) Delete(models ...*gorm_gen_model.AdminDemo) (result gen.ResultInfo, err error) {
	return a.DO.Delete(models)
}

func (a *adminDemoDo) withDO(do gen.Dao) *adminDemoDo {
	a.DO = *do.(*gen.DO)
	return a
}
